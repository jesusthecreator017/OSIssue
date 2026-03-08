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

type Board struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	OwnerUserID *uuid.UUID `json:"owner_user_id"`
	OwnerTeamID *uuid.UUID `json:"owner_team_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type BoardColumn struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Position  int32     `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type BoardStore struct {
	queries *dbsqlc.Queries
}

func (s *BoardStore) CreateBoard(ctx context.Context, board *Board) error {
	params := dbsqlc.CreateBoardParams{
		Name: board.Name,
	}
	if board.OwnerUserID != nil {
		params.OwnerUserID = pgtype.UUID{Bytes: *board.OwnerUserID, Valid: true}
	}
	if board.OwnerTeamID != nil {
		params.OwnerTeamID = pgtype.UUID{Bytes: *board.OwnerTeamID, Valid: true}
	}

	row, err := s.queries.CreateBoard(ctx, params)
	if err != nil {
		return fmt.Errorf("creating board: %w", err)
	}

	board.ID = row.ID
	board.CreatedAt = row.CreatedAt.Time
	board.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (s *BoardStore) GetBoardByID(ctx context.Context, id uuid.UUID) (*Board, error) {
	row, err := s.queries.GetBoardByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting board: %w", err)
	}

	return boardRowToDomain(row), nil
}

func (s *BoardStore) GetPersonalBoard(ctx context.Context, userID uuid.UUID) (*Board, error) {
	row, err := s.queries.GetPersonalBoard(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting personal board: %w", err)
	}

	return boardRowToDomain(row), nil
}

func (s *BoardStore) GetTeamBoard(ctx context.Context, teamID uuid.UUID) (*Board, error) {
	row, err := s.queries.GetTeamBoard(ctx, pgtype.UUID{Bytes: teamID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting team board: %w", err)
	}

	return boardRowToDomain(row), nil
}

func (s *BoardStore) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteBoard(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting board: %w", err)
	}
	return nil
}

func (s *BoardStore) CreateColumn(ctx context.Context, col *BoardColumn) error {
	row, err := s.queries.CreateBoardColumn(ctx, dbsqlc.CreateBoardColumnParams{
		BoardID:  col.BoardID,
		Name:     col.Name,
		Position: col.Position,
	})
	if err != nil {
		return fmt.Errorf("creating board column: %w", err)
	}

	col.ID = row.ID
	col.CreatedAt = row.CreatedAt.Time
	return nil
}

func (s *BoardStore) ListColumns(ctx context.Context, boardID uuid.UUID) ([]*BoardColumn, error) {
	rows, err := s.queries.ListBoardColumns(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("listing board columns: %w", err)
	}

	columns := make([]*BoardColumn, len(rows))
	for i, row := range rows {
		columns[i] = &BoardColumn{
			ID:        row.ID,
			BoardID:   row.BoardID,
			Name:      row.Name,
			Position:  row.Position,
			CreatedAt: row.CreatedAt.Time,
		}
	}
	return columns, nil
}

func (s *BoardStore) UpdateColumn(ctx context.Context, id uuid.UUID, name string) (*BoardColumn, error) {
	row, err := s.queries.UpdateBoardColumn(ctx, dbsqlc.UpdateBoardColumnParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("updating board column: %w", err)
	}

	return &BoardColumn{
		ID:        row.ID,
		BoardID:   row.BoardID,
		Name:      row.Name,
		Position:  row.Position,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (s *BoardStore) ReorderColumn(ctx context.Context, id uuid.UUID, position int32) (*BoardColumn, error) {
	row, err := s.queries.ReorderBoardColumn(ctx, dbsqlc.ReorderBoardColumnParams{
		ID:       id,
		Position: position,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("reordering board column: %w", err)
	}

	return &BoardColumn{
		ID:        row.ID,
		BoardID:   row.BoardID,
		Name:      row.Name,
		Position:  row.Position,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (s *BoardStore) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteBoardColumn(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting board column: %w", err)
	}
	return nil
}

func boardRowToDomain(row dbsqlc.Board) *Board {
	board := &Board{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.OwnerUserID.Valid {
		id := uuid.UUID(row.OwnerUserID.Bytes)
		board.OwnerUserID = &id
	}
	if row.OwnerTeamID.Valid {
		id := uuid.UUID(row.OwnerTeamID.Bytes)
		board.OwnerTeamID = &id
	}
	return board
}
