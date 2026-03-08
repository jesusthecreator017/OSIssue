package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
	"github.com/jesusthecreator017/fswithgo/cmd/api/middleware"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

func (app *application) createIssueHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Title         string  `json:"title"`
		Description   string  `json:"description"`
		Priority      string  `json:"priority"`
		AssigneeID    *string `json:"assignee_id"`
		TeamID        *string `json:"team_id"`
		BoardColumnID *string `json:"board_column_id"`
		DueDate       *string `json:"due_date"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		errs["title"] = "must not be blank"
	} else if len(input.Title) > 255 {
		errs["title"] = "must not be more than 255 characters"
	}

	priority := store.PriorityMedium
	if input.Priority != "" {
		switch store.PriorityType(input.Priority) {
		case store.PriorityLow, store.PriorityMedium, store.PriorityHigh, store.PriorityCritical:
			priority = store.PriorityType(input.Priority)
		default:
			errs["priority"] = "must be one of: Low, Medium, High, Critical"
		}
	}

	var assigneeID *uuid.UUID
	if input.AssigneeID != nil && *input.AssigneeID != "" {
		id, err := uuid.Parse(*input.AssigneeID)
		if err != nil {
			errs["assignee_id"] = "must be a valid UUID"
		} else {
			assigneeID = &id
		}
	}

	var teamID *uuid.UUID
	if input.TeamID != nil && *input.TeamID != "" {
		id, err := uuid.Parse(*input.TeamID)
		if err != nil {
			errs["team_id"] = "must be a valid UUID"
		} else {
			teamID = &id
		}
	}

	var boardColumnID *uuid.UUID
	if input.BoardColumnID != nil && *input.BoardColumnID != "" {
		id, err := uuid.Parse(*input.BoardColumnID)
		if err != nil {
			errs["board_column_id"] = "must be a valid UUID"
		} else {
			boardColumnID = &id
		}
	}

	var dueDate *time.Time
	if input.DueDate != nil && *input.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			errs["due_date"] = "must be a valid RFC3339 date"
		} else {
			dueDate = &t
		}
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	userID := middleware.GetUserID(req)

	issue := &store.Issue{
		Title:         input.Title,
		UserID:        userID,
		Description:   input.Description,
		Priority:      priority,
		AssigneeID:    assigneeID,
		TeamID:        teamID,
		BoardColumnID: boardColumnID,
		DueDate:       dueDate,
	}

	if err := app.store.Issues.Create(req.Context(), issue); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create issue")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"issue": issue})
}

func (app *application) listIssueHandler(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req)

	issueList, err := app.store.Issues.ListByUserID(req.Context(), userID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get issues")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issues": issueList})
}

func (app *application) getIssueHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	issueID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	issue, err := app.store.Issues.GetByID(req.Context(), issueID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "issue not found")
			return
		}
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get issue")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issue": issue})
}

func (app *application) updateIssueHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	issueID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	issue, err := app.store.Issues.GetByID(req.Context(), issueID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "issue not found")
			return
		}
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get issue")
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Priority    *string `json:"priority"`
		AssigneeID  *string `json:"assignee_id"`
		TeamID      *string `json:"team_id"`
		DueDate     *string `json:"due_date"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	if input.Title != nil {
		t := strings.TrimSpace(*input.Title)
		if t == "" {
			errs["title"] = "must not be blank"
		} else if len(t) > 255 {
			errs["title"] = "must not be more than 255 characters"
		} else {
			issue.Title = t
		}
	}

	if input.Description != nil {
		issue.Description = *input.Description
	}

	if input.Priority != nil {
		switch store.PriorityType(*input.Priority) {
		case store.PriorityLow, store.PriorityMedium, store.PriorityHigh, store.PriorityCritical:
			issue.Priority = store.PriorityType(*input.Priority)
		default:
			errs["priority"] = "must be one of: Low, Medium, High, Critical"
		}
	}

	if input.AssigneeID != nil {
		if *input.AssigneeID == "" {
			issue.AssigneeID = nil
		} else {
			aid, err := uuid.Parse(*input.AssigneeID)
			if err != nil {
				errs["assignee_id"] = "must be a valid UUID"
			} else {
				issue.AssigneeID = &aid
			}
		}
	}

	if input.TeamID != nil {
		if *input.TeamID == "" {
			issue.TeamID = nil
		} else {
			tid, err := uuid.Parse(*input.TeamID)
			if err != nil {
				errs["team_id"] = "must be a valid UUID"
			} else {
				issue.TeamID = &tid
			}
		}
	}

	if input.DueDate != nil {
		if *input.DueDate == "" {
			issue.DueDate = nil
		} else {
			t, err := time.Parse(time.RFC3339, *input.DueDate)
			if err != nil {
				errs["due_date"] = "must be a valid RFC3339 date"
			} else {
				issue.DueDate = &t
			}
		}
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	if err := app.store.Issues.Update(req.Context(), issue); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to update issue")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issue": issue})
}

func (app *application) moveIssueHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	issueID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	var input struct {
		BoardColumnID string `json:"board_column_id"`
		Position      int32  `json:"position"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	columnID, err := uuid.Parse(input.BoardColumnID)
	if err != nil {
		errs["board_column_id"] = "must be a valid UUID"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	if err := app.store.Issues.MoveIssue(req.Context(), issueID, columnID, input.Position); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "issue not found")
			return
		}
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to move issue")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "issue moved"})
}

func (app *application) deleteIssueHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	issueID, err := uuid.Parse(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	err = app.store.Issues.Delete(req.Context(), issueID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			helpers.ErrorJson(w, http.StatusNotFound, "issue not found")
			return
		}
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to delete issue")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "issue deleted"})
}
