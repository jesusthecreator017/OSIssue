-- name: CreateTeam :one
INSERT INTO teams (name, description, created_by, avatar_url, max_members)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, created_by, avatar_url, max_members, created_at, updated_at;

-- name: GetTeamByID :one
SELECT t.id, t.name, t.description, t.created_by, t.avatar_url, t.max_members, t.created_at, t.updated_at
FROM teams t
WHERE t.id = $1;

-- name: GetTeamByName :one
SELECT t.id, t.name, t.description, t.created_by, t.avatar_url, t.max_members, t.created_at, t.updated_at
FROM teams t
WHERE t.name = $1;

-- name: ListTeams :many
SELECT t.id, t.name, t.description, t.created_by, t.avatar_url, t.max_members, t.created_at, t.updated_at
FROM teams t
ORDER BY t.created_at DESC;

-- name: AddUserToTeam :one
INSERT INTO team_members (user_id, team_id, role)
VALUES ($1, $2, $3)
RETURNING user_id, team_id, role, joined_at;

-- name: RemoveUserFromTeam :exec
DELETE FROM team_members
WHERE user_id = $1 AND team_id = $2;

-- name: ListTeamMembers :many
SELECT tm.user_id, u.name AS user_name, u.email, tm.role, tm.joined_at
FROM team_members tm
JOIN users u ON u.id = tm.user_id
WHERE tm.team_id = $1
ORDER BY tm.joined_at ASC;

-- name: ListUserTeams :many
SELECT t.id, t.name, t.description, tm.role, tm.joined_at, t.created_at, t.updated_at
FROM team_members tm
JOIN teams t ON t.id = tm.team_id
WHERE tm.user_id = $1
ORDER BY tm.joined_at ASC;

-- name: CountTeamMembers :one
SELECT COUNT(*) AS count
FROM team_members
WHERE team_id = $1;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;
