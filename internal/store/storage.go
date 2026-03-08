package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

var ErrNotFound = errors.New("resource not found")

type Storage struct {
	Issues interface {
		Create(context.Context, *Issue) error
		GetByID(context.Context, uuid.UUID) (*Issue, error)
		List(context.Context) ([]*Issue, error)
		ListByUserID(context.Context, uuid.UUID) ([]*Issue, error)
		ListByTeamID(context.Context, uuid.UUID) ([]*Issue, error)
		ListByBoardID(context.Context, uuid.UUID) ([]*Issue, error)
		Update(context.Context, *Issue) error
		MoveIssue(context.Context, uuid.UUID, uuid.UUID, int32) error
		Delete(context.Context, uuid.UUID) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, uuid.UUID) (*User, error)
		SearchByName(context.Context, string) ([]*User, error)
	}
	Admin interface {
		GetStats(context.Context) (*AdminStats, error)
	}
	Teams interface {
		Create(ctx context.Context, team *Team) error
		GetTeamByID(ctx context.Context, id uuid.UUID) (*Team, error)
		GetTeamByName(ctx context.Context, name string) (*Team, error)
		GetTeamMemberList(ctx context.Context, teamID uuid.UUID) ([]*TeamMember, error)
		GetTeamList(ctx context.Context) ([]*Team, error)
		GetUserTeamsList(ctx context.Context, userID uuid.UUID) ([]*ListUserTeamsRow, error)
		AddUserToTeam(ctx context.Context, teamMember *TeamMember) error
		RemoveUserFromTeam(ctx context.Context, userID uuid.UUID, teamID uuid.UUID) error
		CountMembers(ctx context.Context, teamID uuid.UUID) (int64, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}
	Boards interface {
		CreateBoard(context.Context, *Board) error
		GetBoardByID(context.Context, uuid.UUID) (*Board, error)
		GetPersonalBoard(context.Context, uuid.UUID) (*Board, error)
		GetTeamBoard(context.Context, uuid.UUID) (*Board, error)
		DeleteBoard(context.Context, uuid.UUID) error
		CreateColumn(context.Context, *BoardColumn) error
		ListColumns(context.Context, uuid.UUID) ([]*BoardColumn, error)
		UpdateColumn(context.Context, uuid.UUID, string) (*BoardColumn, error)
		ReorderColumn(context.Context, uuid.UUID, int32) (*BoardColumn, error)
		DeleteColumn(context.Context, uuid.UUID) error
	}
	Labels interface {
		Create(context.Context, *Label) error
		List(context.Context) ([]*Label, error)
		GetByID(context.Context, uuid.UUID) (*Label, error)
		Delete(context.Context, uuid.UUID) error
		AddToIssue(context.Context, uuid.UUID, uuid.UUID) error
		RemoveFromIssue(context.Context, uuid.UUID, uuid.UUID) error
		ListForIssue(context.Context, uuid.UUID) ([]*Label, error)
	}
}

func NewStorage(pool *pgxpool.Pool) Storage {
	queries := dbsqlc.New(pool)

	return Storage{
		Issues: &IssueStore{queries: queries},
		Users:  &UserStore{queries: queries},
		Admin:  &AdminStore{queries: queries},
		Teams:  &TeamStore{queries: queries},
		Boards: &BoardStore{queries: queries},
		Labels: &LabelStore{queries: queries},
	}
}
