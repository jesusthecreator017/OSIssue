-- name: CountUsers :one
SELECT COUNT(*) AS total_users FROM users;

-- name: CountIssuesByPriority :many
SELECT priority, COUNT(*) AS count FROM issues GROUP BY priority;
