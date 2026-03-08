package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
	"github.com/jesusthecreator017/fswithgo/cmd/api/middleware"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

func (app *application) getPersonalBoardHandler(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req)

	board, err := app.store.Boards.GetPersonalBoard(req.Context(), userID)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get board")
			return
		}

		// Auto-create personal board with default columns
		board = &store.Board{
			Name:        "My Board",
			OwnerUserID: &userID,
		}
		if err := app.store.Boards.CreateBoard(req.Context(), board); err != nil {
			helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create board")
			return
		}

		defaultColumns := []string{"Todo", "In Progress", "Done"}
		for i, name := range defaultColumns {
			col := &store.BoardColumn{
				BoardID:  board.ID,
				Name:     name,
				Position: int32(i),
			}
			if err := app.store.Boards.CreateColumn(req.Context(), col); err != nil {
				helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create default columns")
				return
			}
		}
	}

	columns, err := app.store.Boards.ListColumns(req.Context(), board.ID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to list columns")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{
		"board":   board,
		"columns": columns,
	})
}

func (app *application) getTeamBoardHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	teamID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	board, err := app.store.Boards.GetTeamBoard(req.Context(), teamID)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get board")
			return
		}

		// Auto-create team board with default columns
		board = &store.Board{
			Name:        "Team Board",
			OwnerTeamID: &teamID,
		}
		if err := app.store.Boards.CreateBoard(req.Context(), board); err != nil {
			helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create board")
			return
		}

		defaultColumns := []string{"Todo", "In Progress", "Done"}
		for i, name := range defaultColumns {
			col := &store.BoardColumn{
				BoardID:  board.ID,
				Name:     name,
				Position: int32(i),
			}
			if err := app.store.Boards.CreateColumn(req.Context(), col); err != nil {
				helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create default columns")
				return
			}
		}
	}

	columns, err := app.store.Boards.ListColumns(req.Context(), board.ID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to list columns")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{
		"board":   board,
		"columns": columns,
	})
}

func (app *application) listBoardIssuesHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	boardID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	issues, err := app.store.Issues.ListByBoardID(req.Context(), boardID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to list issues")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issues": issues})
}

func (app *application) createBoardColumnHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	boardID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	var input struct {
		Name     string `json:"name"`
		Position int32  `json:"position"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		errs["name"] = "must not be blank"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	col := &store.BoardColumn{
		BoardID:  boardID,
		Name:     input.Name,
		Position: input.Position,
	}

	if err := app.store.Boards.CreateColumn(req.Context(), col); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create column")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"column": col})
}

func (app *application) updateBoardColumnHandler(w http.ResponseWriter, req *http.Request) {
	colID, err := uuid.Parse(req.PathValue("colId"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "colId must be a valid UUID")
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		errs["name"] = "must not be blank"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	col, err := app.store.Boards.UpdateColumn(req.Context(), colID, input.Name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "column not found")
			return
		}
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to update column")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"column": col})
}

func (app *application) reorderBoardColumnHandler(w http.ResponseWriter, req *http.Request) {
	colID, err := uuid.Parse(req.PathValue("colId"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "colId must be a valid UUID")
		return
	}

	var input struct {
		Position int32 `json:"position"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	col, err := app.store.Boards.ReorderColumn(req.Context(), colID, input.Position)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "column not found")
			return
		}
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to reorder column")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"column": col})
}

func (app *application) deleteBoardColumnHandler(w http.ResponseWriter, req *http.Request) {
	colID, err := uuid.Parse(req.PathValue("colId"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "colId must be a valid UUID")
		return
	}

	if err := app.store.Boards.DeleteColumn(req.Context(), colID); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to delete column")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "column deleted"})
}
