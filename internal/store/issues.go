package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

type PriorityType string

const (
	PriorityLow      PriorityType = "Low"
	PriorityMedium   PriorityType = "Medium"
	PriorityHigh     PriorityType = "High"
	PriorityCritical PriorityType = "Critical"
)

type Issue struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"user_id"`
	UserName        string       `json:"user_name"`
	AssigneeID      *uuid.UUID   `json:"assignee_id"`
	AssigneeName    string       `json:"assignee_name"`
	TeamID          *uuid.UUID   `json:"team_id"`
	BoardColumnID   *uuid.UUID   `json:"board_column_id"`
	BoardColumnName string       `json:"board_column_name"`
	Position        int32        `json:"position"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Priority        PriorityType `json:"priority"`
	DueDate         *time.Time   `json:"due_date"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type IssueStore struct {
	queries *dbsqlc.Queries
}

func (s *IssueStore) Create(ctx context.Context, issue *Issue) error {
	params := dbsqlc.CreateIssueParams{
		Title:       issue.Title,
		UserID:      issue.UserID,
		Description: issue.Description,
		Priority:    string(issue.Priority),
		Position:    issue.Position,
	}
	if issue.AssigneeID != nil {
		params.AssigneeID = pgtype.UUID{Bytes: *issue.AssigneeID, Valid: true}
	}
	if issue.TeamID != nil {
		params.TeamID = pgtype.UUID{Bytes: *issue.TeamID, Valid: true}
	}
	if issue.BoardColumnID != nil {
		params.BoardColumnID = pgtype.UUID{Bytes: *issue.BoardColumnID, Valid: true}
	}
	if issue.DueDate != nil {
		params.DueDate = pgtype.Timestamptz{Time: *issue.DueDate, Valid: true}
	}

	row, err := s.queries.CreateIssue(ctx, params)
	if err != nil {
		return fmt.Errorf("creating issue: %w", err)
	}

	issue.ID = row.ID
	issue.CreatedAt = row.CreatedAt.Time
	issue.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (s *IssueStore) GetByID(ctx context.Context, id uuid.UUID) (*Issue, error) {
	row, err := s.queries.GetIssueByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting issue: %w", err)
	}

	return issueRowToDomain(row), nil
}

func (s *IssueStore) List(ctx context.Context) ([]*Issue, error) {
	rows, err := s.queries.ListIssues(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing issues: %w", err)
	}

	issues := make([]*Issue, len(rows))
	for i, row := range rows {
		issues[i] = listRowToDomain(row)
	}
	return issues, nil
}

func (s *IssueStore) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Issue, error) {
	rows, err := s.queries.ListIssuesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing issues by user: %w", err)
	}

	issues := make([]*Issue, len(rows))
	for i, row := range rows {
		issues[i] = &Issue{
			ID:              row.ID,
			UserID:          row.UserID,
			UserName:        row.UserName,
			AssigneeName:    row.AssigneeName,
			BoardColumnName: row.BoardColumnName,
			Position:        row.Position,
			Title:           row.Title,
			Description:     row.Description,
			Priority:        PriorityType(row.Priority),
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
		}
		issues[i].AssigneeID = pgtypeUUIDToPtr(row.AssigneeID)
		issues[i].TeamID = pgtypeUUIDToPtr(row.TeamID)
		issues[i].BoardColumnID = pgtypeUUIDToPtr(row.BoardColumnID)
		if row.DueDate.Valid {
			t := row.DueDate.Time
			issues[i].DueDate = &t
		}
	}
	return issues, nil
}

func (s *IssueStore) ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]*Issue, error) {
	rows, err := s.queries.ListIssuesByTeamID(ctx, pgtype.UUID{Bytes: teamID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("listing issues by team: %w", err)
	}

	issues := make([]*Issue, len(rows))
	for i, row := range rows {
		issues[i] = teamIssueRowToDomain(row)
	}
	return issues, nil
}

func (s *IssueStore) ListByBoardID(ctx context.Context, boardID uuid.UUID) ([]*Issue, error) {
	rows, err := s.queries.ListIssuesByBoardID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("listing issues by board: %w", err)
	}

	issues := make([]*Issue, len(rows))
	for i, row := range rows {
		issues[i] = boardIssueRowToDomain(row)
	}
	return issues, nil
}

func (s *IssueStore) Update(ctx context.Context, issue *Issue) error {
	params := dbsqlc.UpdateIssueParams{
		ID:          issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Priority:    string(issue.Priority),
	}
	if issue.AssigneeID != nil {
		params.AssigneeID = pgtype.UUID{Bytes: *issue.AssigneeID, Valid: true}
	}
	if issue.TeamID != nil {
		params.TeamID = pgtype.UUID{Bytes: *issue.TeamID, Valid: true}
	}
	if issue.DueDate != nil {
		params.DueDate = pgtype.Timestamptz{Time: *issue.DueDate, Valid: true}
	}

	row, err := s.queries.UpdateIssue(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("updating issue: %w", err)
	}

	issue.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (s *IssueStore) MoveIssue(ctx context.Context, id uuid.UUID, boardColumnID uuid.UUID, position int32) error {
	_, err := s.queries.MoveIssue(ctx, dbsqlc.MoveIssueParams{
		ID:            id,
		BoardColumnID: pgtype.UUID{Bytes: boardColumnID, Valid: true},
		Position:      position,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("moving issue: %w", err)
	}
	return nil
}

func (s *IssueStore) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteIssue(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("deleting issue: %w", err)
	}
	return nil
}

func pgtypeUUIDToPtr(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	id := uuid.UUID(u.Bytes)
	return &id
}

func pgtypeTsToPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

func issueRowToDomain(row dbsqlc.GetIssueByIDRow) *Issue {
	return &Issue{
		ID:              row.ID,
		UserID:          row.UserID,
		UserName:        row.UserName,
		AssigneeID:      pgtypeUUIDToPtr(row.AssigneeID),
		AssigneeName:    row.AssigneeName,
		TeamID:          pgtypeUUIDToPtr(row.TeamID),
		BoardColumnID:   pgtypeUUIDToPtr(row.BoardColumnID),
		BoardColumnName: row.BoardColumnName,
		Position:        row.Position,
		Title:           row.Title,
		Description:     row.Description,
		Priority:        PriorityType(row.Priority),
		DueDate:         pgtypeTsToPtr(row.DueDate),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

func listRowToDomain(row dbsqlc.ListIssuesRow) *Issue {
	return &Issue{
		ID:              row.ID,
		UserID:          row.UserID,
		UserName:        row.UserName,
		AssigneeID:      pgtypeUUIDToPtr(row.AssigneeID),
		AssigneeName:    row.AssigneeName,
		TeamID:          pgtypeUUIDToPtr(row.TeamID),
		BoardColumnID:   pgtypeUUIDToPtr(row.BoardColumnID),
		BoardColumnName: row.BoardColumnName,
		Position:        row.Position,
		Title:           row.Title,
		Description:     row.Description,
		Priority:        PriorityType(row.Priority),
		DueDate:         pgtypeTsToPtr(row.DueDate),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

func teamIssueRowToDomain(row dbsqlc.ListIssuesByTeamIDRow) *Issue {
	return &Issue{
		ID:              row.ID,
		UserID:          row.UserID,
		UserName:        row.UserName,
		AssigneeID:      pgtypeUUIDToPtr(row.AssigneeID),
		AssigneeName:    row.AssigneeName,
		TeamID:          pgtypeUUIDToPtr(row.TeamID),
		BoardColumnID:   pgtypeUUIDToPtr(row.BoardColumnID),
		BoardColumnName: row.BoardColumnName,
		Position:        row.Position,
		Title:           row.Title,
		Description:     row.Description,
		Priority:        PriorityType(row.Priority),
		DueDate:         pgtypeTsToPtr(row.DueDate),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

func boardIssueRowToDomain(row dbsqlc.ListIssuesByBoardIDRow) *Issue {
	return &Issue{
		ID:              row.ID,
		UserID:          row.UserID,
		UserName:        row.UserName,
		AssigneeID:      pgtypeUUIDToPtr(row.AssigneeID),
		AssigneeName:    row.AssigneeName,
		TeamID:          pgtypeUUIDToPtr(row.TeamID),
		BoardColumnID:   pgtypeUUIDToPtr(row.BoardColumnID),
		BoardColumnName: row.BoardColumnName,
		Position:        row.Position,
		Title:           row.Title,
		Description:     row.Description,
		Priority:        PriorityType(row.Priority),
		DueDate:         pgtypeTsToPtr(row.DueDate),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}
