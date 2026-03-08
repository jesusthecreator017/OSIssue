-- name: CreateBoardColumn :one
INSERT INTO board_columns (board_id, name, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListBoardColumns :many
SELECT * FROM board_columns
WHERE board_id = $1
ORDER BY position;

-- name: UpdateBoardColumn :one
UPDATE board_columns
SET name = $2
WHERE id = $1
RETURNING *;

-- name: ReorderBoardColumn :one
UPDATE board_columns
SET position = $2
WHERE id = $1
RETURNING *;

-- name: DeleteBoardColumn :exec
DELETE FROM board_columns WHERE id = $1;
