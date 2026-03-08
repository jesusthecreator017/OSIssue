-- name: CreateLabel :one
INSERT INTO labels (name, color)
VALUES ($1, $2)
RETURNING *;

-- name: ListLabels :many
SELECT * FROM labels ORDER BY name;

-- name: GetLabelByID :one
SELECT * FROM labels WHERE id = $1;

-- name: DeleteLabel :exec
DELETE FROM labels WHERE id = $1;

-- name: AddLabelToIssue :exec
INSERT INTO issue_labels (issue_id, label_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveLabelFromIssue :exec
DELETE FROM issue_labels
WHERE issue_id = $1 AND label_id = $2;

-- name: ListLabelsForIssue :many
SELECT l.*
FROM labels l
JOIN issue_labels il ON il.label_id = l.id
WHERE il.issue_id = $1
ORDER BY l.name;
