package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jesusthecreator017/fswithgo/cmd/api/middleware"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr       string
	db         dbConfig
	corsOrigin string
	jwtSecret  string
}

type dbConfig struct {
	dsn         string
	maxConns    int
	maxIdleTime time.Duration
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	// Issues
	mux.Handle("GET /v1/issues", middleware.RequiredAuth(http.HandlerFunc(app.listIssueHandler)))
	mux.HandleFunc("GET /v1/issues/{id}", app.getIssueHandler)
	mux.Handle("POST /v1/issues", middleware.RequiredAuth(http.HandlerFunc(app.createIssueHandler)))
	mux.Handle("PATCH /v1/issues/{id}", middleware.RequiredAuth(http.HandlerFunc(app.updateIssueHandler)))
	mux.Handle("PATCH /v1/issues/{id}/move", middleware.RequiredAuth(http.HandlerFunc(app.moveIssueHandler)))
	mux.Handle("DELETE /v1/issues/{id}", middleware.RequiredAuth(http.HandlerFunc(app.deleteIssueHandler)))

	// Boards
	mux.Handle("GET /v1/me/board", middleware.RequiredAuth(http.HandlerFunc(app.getPersonalBoardHandler)))
	mux.Handle("GET /v1/teams/{id}/board", middleware.RequiredAuth(http.HandlerFunc(app.getTeamBoardHandler)))
	mux.Handle("GET /v1/boards/{id}/issues", middleware.RequiredAuth(http.HandlerFunc(app.listBoardIssuesHandler)))

	// Board columns
	mux.Handle("POST /v1/boards/{id}/columns", middleware.RequiredAuth(http.HandlerFunc(app.createBoardColumnHandler)))
	mux.Handle("PATCH /v1/boards/{id}/columns/{colId}", middleware.RequiredAuth(http.HandlerFunc(app.updateBoardColumnHandler)))
	mux.Handle("PATCH /v1/boards/{id}/columns/{colId}/reorder", middleware.RequiredAuth(http.HandlerFunc(app.reorderBoardColumnHandler)))
	mux.Handle("DELETE /v1/boards/{id}/columns/{colId}", middleware.RequiredAuth(http.HandlerFunc(app.deleteBoardColumnHandler)))

	// Labels
	mux.Handle("GET /v1/labels", middleware.RequiredAuth(http.HandlerFunc(app.listLabelsHandler)))
	mux.Handle("POST /v1/labels", middleware.RequiredAuth(http.HandlerFunc(app.createLabelHandler)))
	mux.Handle("POST /v1/issues/{id}/labels", middleware.RequiredAuth(http.HandlerFunc(app.addLabelToIssueHandler)))
	mux.Handle("DELETE /v1/issues/{id}/labels/{labelId}", middleware.RequiredAuth(http.HandlerFunc(app.removeLabelFromIssueHandler)))
	mux.Handle("GET /v1/issues/{id}/labels", middleware.RequiredAuth(http.HandlerFunc(app.listIssueLabelsHandler)))

	// Users
	mux.HandleFunc("POST /v1/users/register", app.registerUserHandler)
	mux.HandleFunc("POST /v1/users/login", app.loginUserHandler)
	mux.Handle("GET /v1/users/search", middleware.RequiredAuth(http.HandlerFunc(app.searchUsersHandler)))

	// Teams
	mux.Handle("POST /v1/teams", middleware.RequiredAuth(http.HandlerFunc(app.createTeamHandler)))
	mux.Handle("GET /v1/teams", middleware.RequiredAuth(http.HandlerFunc(app.listTeamsHandler)))
	mux.Handle("GET /v1/teams/{id}", middleware.RequiredAuth(http.HandlerFunc(app.getTeamHandler)))
	mux.Handle("GET /v1/teams/{id}/members", middleware.RequiredAuth(http.HandlerFunc(app.getTeamMembersHandler)))
	mux.Handle("POST /v1/teams/{id}/members", middleware.RequiredAuth(http.HandlerFunc(app.addTeamMemberHandler)))
	mux.Handle("DELETE /v1/teams/{id}/members/{userId}", middleware.RequiredAuth(http.HandlerFunc(app.removeTeamMemberHandler)))
	mux.Handle("DELETE /v1/teams/{id}", middleware.RequiredAuth(http.HandlerFunc(app.deleteTeamHandler)))
	mux.Handle("GET /v1/me/teams", middleware.RequiredAuth(http.HandlerFunc(app.getUserTeamsHandler)))

	// Admin
	mux.Handle("GET /v1/admin/stats", middleware.RequiredAuth(middleware.RequiredAdmin(http.HandlerFunc(app.adminStatsHandler))))

	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	stack := middleware.CreateStack(
		middleware.CORS(app.config.corsOrigin),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logging,
		middleware.GlobalAuth(app.config.jwtSecret),
		middleware.Timeout(time.Second*25),
	)

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      stack(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("Listening on port%s\n", app.config.addr)
	return srv.ListenAndServe()
}
