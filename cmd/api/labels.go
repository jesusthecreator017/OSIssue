package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

func (app *application) createLabelHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Color string `json:"color"`
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
	if input.Color == "" {
		input.Color = "#6b7280"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	label := &store.Label{
		Name:  input.Name,
		Color: input.Color,
	}

	if err := app.store.Labels.Create(req.Context(), label); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create label")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"label": label})
}

func (app *application) listLabelsHandler(w http.ResponseWriter, req *http.Request) {
	labels, err := app.store.Labels.List(req.Context())
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to list labels")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"labels": labels})
}

func (app *application) addLabelToIssueHandler(w http.ResponseWriter, req *http.Request) {
	issueID, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	var input struct {
		LabelID string `json:"label_id"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	labelID, err := uuid.Parse(input.LabelID)
	if err != nil {
		helpers.ValidationErrorJson(w, map[string]string{"label_id": "must be a valid UUID"})
		return
	}

	if err := app.store.Labels.AddToIssue(req.Context(), issueID, labelID); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to add label to issue")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "label added"})
}

func (app *application) removeLabelFromIssueHandler(w http.ResponseWriter, req *http.Request) {
	issueID, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	labelID, err := uuid.Parse(req.PathValue("labelId"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "labelId must be a valid UUID")
		return
	}

	if err := app.store.Labels.RemoveFromIssue(req.Context(), issueID, labelID); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to remove label from issue")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "label removed"})
}

func (app *application) listIssueLabelsHandler(w http.ResponseWriter, req *http.Request) {
	issueID, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	labels, err := app.store.Labels.ListForIssue(req.Context(), issueID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to list labels")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"labels": labels})
}
