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

func (app *application) createTeamHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
		MaxMembers  int32  `json:"max_members"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		errs["name"] = "must not be blank"
	} else if len(input.Name) > 127 {
		errs["name"] = "must not be more than 127 characters"
	}

	input.Description = strings.TrimSpace(input.Description)
	if len(input.Description) > 500 {
		errs["description"] = "must not be more than 500 characters"
	}

	if input.MaxMembers < 0 {
		errs["max_members"] = "must not be negative"
	}
	if input.MaxMembers == 0 {
		input.MaxMembers = 50
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	userID := middleware.GetUserID(req)

	team := &store.Team{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   userID,
		AvatarURL:   input.AvatarURL,
		MaxMembers:  input.MaxMembers,
	}

	if err := app.store.Teams.Create(req.Context(), team); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create team")
		return
	}

	// Auto-add creator as owner
	owner := &store.TeamMember{
		UserID: userID,
		TeamID: team.ID,
		Role:   store.TeamRoleOwner,
	}
	if err := app.store.Teams.AddUserToTeam(req.Context(), owner); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to add creator as owner")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"team": team})
}

func (app *application) getTeamHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	// Check that an ID is passed in
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be provided")
		return
	}

	// Check that the id is a uuid
	teamID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a uuid")
		return
	}

	// Get the team with the given id
	team, err := app.store.Teams.GetTeamByID(req.Context(), teamID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "team not found")
		} else {
			helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get team")
		}
		return
	}

	// Return the team in the response
	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"team": team})
}

func (app *application) listTeamsHandler(w http.ResponseWriter, req *http.Request) {
	teamList, err := app.store.Teams.GetTeamList(req.Context())
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get team list")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"teams": teamList})
}

func (app *application) getUserTeamsHandler(w http.ResponseWriter, req *http.Request) {
	userId := middleware.GetUserID(req)

	userTeams, err := app.store.Teams.GetUserTeamsList(req.Context(), userId)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get user teams list")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"teams": userTeams})
}

func (app *application) getTeamMembersHandler(w http.ResponseWriter, req *http.Request) {
	teamIdStr := req.PathValue("id")

	// Check that an ID is passed in
	if teamIdStr == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "teamid must not be empty")
		return
	}

	// Check that the id is a uuid
	teamID, err := uuid.Parse(teamIdStr)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "teamid must be a uuid")
		return
	}

	teamMembers, err := app.store.Teams.GetTeamMemberList(req.Context(), teamID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get team members list")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"members": teamMembers})
}

func (app *application) addTeamMemberHandler(w http.ResponseWriter, req *http.Request) {
	// Get team ID from URL path
	teamIDStr := req.PathValue("id")
	if teamIDStr == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "team id must be provided")
		return
	}
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "team id must be a uuid")
		return
	}

	var input struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	input.UserID = strings.TrimSpace(input.UserID)
	if input.UserID == "" {
		errs["user_id"] = "must not be blank"
	} else if _, err := uuid.Parse(input.UserID); err != nil {
		errs["user_id"] = "must be a valid uuid"
	}

	input.Role = strings.TrimSpace(input.Role)
	if input.Role == "" {
		errs["role"] = "must not be blank"
	} else if input.Role != string(store.TeamRoleMember) && input.Role != string(store.TeamRoleAdmin) {
		errs["role"] = "must be either 'member' or 'admin'"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	// Enforce max_members limit
	team, err := app.store.Teams.GetTeamByID(req.Context(), teamID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get team")
		return
	}

	count, err := app.store.Teams.CountMembers(req.Context(), teamID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to count members")
		return
	}

	if count >= int64(team.MaxMembers) {
		helpers.ErrorJson(w, http.StatusConflict, "team has reached its maximum number of members")
		return
	}

	userID, _ := uuid.Parse(input.UserID)

	teamMember := &store.TeamMember{
		UserID: userID,
		TeamID: teamID,
		Role:   store.TeamRole(input.Role),
	}

	if err := app.store.Teams.AddUserToTeam(req.Context(), teamMember); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to add team member")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"team_member": teamMember})
}

func (app *application) removeTeamMemberHandler(w http.ResponseWriter, req *http.Request) {
	// Get Team ID from url path
	teamIdStr := req.PathValue("id")
	if teamIdStr == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "team id must be provided")
		return
	}
	teamID, err := uuid.Parse(teamIdStr)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "team id must be a uuid")
		return
	}

	// Get User ID from url path
	userIdStr := req.PathValue("userId")
	if userIdStr == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "user id must be provided")
		return
	}
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "user id must be a uuid")
		return
	}

	if err := app.store.Teams.RemoveUserFromTeam(req.Context(), userID, teamID); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to remove team member")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "team member removed"})
}

func (app *application) deleteTeamHandler(w http.ResponseWriter, req *http.Request) {
	teamIdStr := req.PathValue("id")
	if teamIdStr == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "team id must be provided")
		return
	}
	teamID, err := uuid.Parse(teamIdStr)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "team id must be a uuid")
		return
	}

	// Only the team creator can delete
	team, err := app.store.Teams.GetTeamByID(req.Context(), teamID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "team not found")
		} else {
			helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get team")
		}
		return
	}

	userID := middleware.GetUserID(req)
	if team.CreatedBy != userID {
		helpers.ErrorJson(w, http.StatusForbidden, "only the team creator can delete this team")
		return
	}

	if err := app.store.Teams.Delete(req.Context(), teamID); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to delete team")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "team deleted"})
}
