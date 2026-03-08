-- name: CreateIssue :one
INSERT INTO issues (title, user_id, description, priority, assignee_id, team_id, board_column_id, position, due_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetIssueByID :one
SELECT i.*,
       u.name AS user_name,
       COALESCE(a.name, '') AS assignee_name,
       COALESCE(bc.name, '') AS board_column_name
FROM issues i
JOIN users u ON u.id = i.user_id
LEFT JOIN users a ON a.id = i.assignee_id
LEFT JOIN board_columns bc ON bc.id = i.board_column_id
WHERE i.id = $1;

-- name: ListIssues :many
SELECT i.*,
       u.name AS user_name,
       COALESCE(a.name, '') AS assignee_name,
       COALESCE(bc.name, '') AS board_column_name
FROM issues i
JOIN users u ON u.id = i.user_id
LEFT JOIN users a ON a.id = i.assignee_id
LEFT JOIN board_columns bc ON bc.id = i.board_column_id
ORDER BY i.created_at DESC;

-- name: ListIssuesByUserID :many
SELECT i.*,
       u.name AS user_name,
       COALESCE(a.name, '') AS assignee_name,
       COALESCE(bc.name, '') AS board_column_name
FROM issues i
JOIN users u ON u.id = i.user_id
LEFT JOIN users a ON a.id = i.assignee_id
LEFT JOIN board_columns bc ON bc.id = i.board_column_id
WHERE i.user_id = $1
ORDER BY i.created_at DESC;

-- name: ListIssuesByTeamID :many
SELECT i.*,
       u.name AS user_name,
       COALESCE(a.name, '') AS assignee_name,
       COALESCE(bc.name, '') AS board_column_name
FROM issues i
JOIN users u ON u.id = i.user_id
LEFT JOIN users a ON a.id = i.assignee_id
LEFT JOIN board_columns bc ON bc.id = i.board_column_id
WHERE i.team_id = $1
ORDER BY i.created_at DESC;

-- name: ListIssuesByBoardColumnID :many
SELECT i.*,
       u.name AS user_name,
       COALESCE(a.name, '') AS assignee_name,
       COALESCE(bc.name, '') AS board_column_name
FROM issues i
JOIN users u ON u.id = i.user_id
LEFT JOIN users a ON a.id = i.assignee_id
LEFT JOIN board_columns bc ON bc.id = i.board_column_id
WHERE i.board_column_id = $1
ORDER BY i.position;

-- name: ListIssuesByBoardID :many
SELECT i.*,
       u.name AS user_name,
       COALESCE(a.name, '') AS assignee_name,
       COALESCE(bc.name, '') AS board_column_name
FROM issues i
JOIN board_columns bc ON bc.id = i.board_column_id
JOIN users u ON u.id = i.user_id
LEFT JOIN users a ON a.id = i.assignee_id
WHERE bc.board_id = $1
ORDER BY bc.position, i.position;

-- name: UpdateIssue :one
UPDATE issues
SET title = $2,
    description = $3,
    priority = $4,
    assignee_id = $5,
    team_id = $6,
    due_date = $7,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: MoveIssue :one
UPDATE issues
SET board_column_id = $2,
    position = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteIssue :exec
DELETE FROM issues
WHERE id = $1;
