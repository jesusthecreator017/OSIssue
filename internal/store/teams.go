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

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

type Team struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`
	AvatarURL   string    `json:"avatar_url"`
	MaxMembers  int32     `json:"max_members"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeamMember struct {
	UserID   uuid.UUID `json:"user_id"`
	TeamID   uuid.UUID `json:"team_id"`
	UserName string    `json:"user_name"`
	Email    string    `json:"email"`
	Role     TeamRole  `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type ListUserTeamsRow struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Role        TeamRole  `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeamStore struct {
	queries *dbsqlc.Queries
}

func (t *TeamStore) Create(ctx context.Context, team *Team) error {
	row, err := t.queries.CreateTeam(ctx, dbsqlc.CreateTeamParams{
		Name:        team.Name,
		Description: team.Description,
		CreatedBy:   team.CreatedBy,
		AvatarUrl:   team.AvatarURL,
		MaxMembers:  team.MaxMembers,
	})
	if err != nil {
		return fmt.Errorf("creating team: %w", err)
	}

	team.ID = row.ID
	team.CreatedAt = row.CreatedAt.Time
	team.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (t *TeamStore) GetTeamByID(ctx context.Context, id uuid.UUID) (*Team, error) {
	row, err := t.queries.GetTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting team: %w", err)
	}

	return teamRowToDomain(row.ID, row.Name, row.Description, row.CreatedBy, row.AvatarUrl, row.MaxMembers, row.CreatedAt.Time, row.UpdatedAt.Time), nil
}

func (t *TeamStore) GetTeamByName(ctx context.Context, name string) (*Team, error) {
	row, err := t.queries.GetTeamByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting team: %w", err)
	}

	return teamRowToDomain(row.ID, row.Name, row.Description, row.CreatedBy, row.AvatarUrl, row.MaxMembers, row.CreatedAt.Time, row.UpdatedAt.Time), nil
}

func (t *TeamStore) GetTeamMemberList(ctx context.Context, teamID uuid.UUID) ([]*TeamMember, error) {
	rows, err := t.queries.ListTeamMembers(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("getting team members: %w", err)
	}

	teamMembers := make([]*TeamMember, len(rows))
	for i, row := range rows {
		teamMembers[i] = &TeamMember{
			UserID:   row.UserID,
			TeamID:   teamID,
			UserName: row.UserName,
			Email:    row.Email,
			Role:     TeamRole(row.Role),
			JoinedAt: row.JoinedAt.Time,
		}
	}
	return teamMembers, nil
}

func (t *TeamStore) GetTeamList(ctx context.Context) ([]*Team, error) {
	rows, err := t.queries.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team list: %w", err)
	}

	teams := make([]*Team, len(rows))
	for i, row := range rows {
		teams[i] = teamRowToDomain(row.ID, row.Name, row.Description, row.CreatedBy, row.AvatarUrl, row.MaxMembers, row.CreatedAt.Time, row.UpdatedAt.Time)
	}
	return teams, nil
}

func (t *TeamStore) GetUserTeamsList(ctx context.Context, userID uuid.UUID) ([]*ListUserTeamsRow, error) {
	rows, err := t.queries.ListUserTeams(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user teams list: %w", err)
	}

	teams := make([]*ListUserTeamsRow, len(rows))
	for i, row := range rows {
		teams[i] = &ListUserTeamsRow{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Role:        TeamRole(row.Role),
			JoinedAt:    row.JoinedAt.Time,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
	}
	return teams, nil
}

// Mutations

func (t *TeamStore) AddUserToTeam(ctx context.Context, teamMember *TeamMember) error {
	row, err := t.queries.AddUserToTeam(ctx, dbsqlc.AddUserToTeamParams{
		UserID: teamMember.UserID,
		TeamID: teamMember.TeamID,
		Role:   string(teamMember.Role),
	})
	if err != nil {
		return fmt.Errorf("adding user to team: %w", err)
	}

	teamMember.JoinedAt = row.JoinedAt.Time
	return nil
}

func (t *TeamStore) RemoveUserFromTeam(ctx context.Context, userID uuid.UUID, teamID uuid.UUID) error {
	err := t.queries.RemoveUserFromTeam(ctx, dbsqlc.RemoveUserFromTeamParams{
		UserID: userID,
		TeamID: teamID,
	})
	if err != nil {
		return fmt.Errorf("removing user from team: %w", err)
	}
	return nil
}

func (t *TeamStore) CountMembers(ctx context.Context, teamID uuid.UUID) (int64, error) {
	count, err := t.queries.CountTeamMembers(ctx, teamID)
	if err != nil {
		return 0, fmt.Errorf("counting team members: %w", err)
	}
	return count, nil
}

func (t *TeamStore) Delete(ctx context.Context, id uuid.UUID) error {
	err := t.queries.DeleteTeam(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting team: %w", err)
	}
	return nil
}

func teamRowToDomain(id uuid.UUID, name, description string, createdBy uuid.UUID, avatarURL string, maxMembers int32, createdAt, updatedAt time.Time) *Team {
	return &Team{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		AvatarURL:   avatarURL,
		MaxMembers:  maxMembers,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
