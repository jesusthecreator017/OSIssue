package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

type Label struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type LabelStore struct {
	queries *dbsqlc.Queries
}

func (s *LabelStore) Create(ctx context.Context, label *Label) error {
	row, err := s.queries.CreateLabel(ctx, dbsqlc.CreateLabelParams{
		Name:  label.Name,
		Color: label.Color,
	})
	if err != nil {
		return fmt.Errorf("creating label: %w", err)
	}

	label.ID = row.ID
	label.CreatedAt = row.CreatedAt.Time
	return nil
}

func (s *LabelStore) List(ctx context.Context) ([]*Label, error) {
	rows, err := s.queries.ListLabels(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing labels: %w", err)
	}

	labels := make([]*Label, len(rows))
	for i, row := range rows {
		labels[i] = &Label{
			ID:        row.ID,
			Name:      row.Name,
			Color:     row.Color,
			CreatedAt: row.CreatedAt.Time,
		}
	}
	return labels, nil
}

func (s *LabelStore) GetByID(ctx context.Context, id uuid.UUID) (*Label, error) {
	row, err := s.queries.GetLabelByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting label: %w", err)
	}

	return &Label{
		ID:        row.ID,
		Name:      row.Name,
		Color:     row.Color,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (s *LabelStore) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteLabel(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting label: %w", err)
	}
	return nil
}

func (s *LabelStore) AddToIssue(ctx context.Context, issueID, labelID uuid.UUID) error {
	err := s.queries.AddLabelToIssue(ctx, dbsqlc.AddLabelToIssueParams{
		IssueID: issueID,
		LabelID: labelID,
	})
	if err != nil {
		return fmt.Errorf("adding label to issue: %w", err)
	}
	return nil
}

func (s *LabelStore) RemoveFromIssue(ctx context.Context, issueID, labelID uuid.UUID) error {
	err := s.queries.RemoveLabelFromIssue(ctx, dbsqlc.RemoveLabelFromIssueParams{
		IssueID: issueID,
		LabelID: labelID,
	})
	if err != nil {
		return fmt.Errorf("removing label from issue: %w", err)
	}
	return nil
}

func (s *LabelStore) ListForIssue(ctx context.Context, issueID uuid.UUID) ([]*Label, error) {
	rows, err := s.queries.ListLabelsForIssue(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("listing labels for issue: %w", err)
	}

	labels := make([]*Label, len(rows))
	for i, row := range rows {
		labels[i] = &Label{
			ID:        row.ID,
			Name:      row.Name,
			Color:     row.Color,
			CreatedAt: row.CreatedAt.Time,
		}
	}
	return labels, nil
}
