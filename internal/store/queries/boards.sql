-- name: CreateBoard :one
INSERT INTO boards (name, owner_user_id, owner_team_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoardByID :one
SELECT * FROM boards WHERE id = $1;

-- name: GetPersonalBoard :one
SELECT * FROM boards WHERE owner_user_id = $1;

-- name: GetTeamBoard :one
SELECT * FROM boards WHERE owner_team_id = $1;

-- name: DeleteBoard :exec
DELETE FROM boards WHERE id = $1;
