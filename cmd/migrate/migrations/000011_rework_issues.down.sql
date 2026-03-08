-- 000011_rework_issues.down.sql
-- Revert issues rework

-- 1. Recreate status_type enum
CREATE TYPE status_type AS ENUM ('Incomplete', 'In-Progress', 'Complete');

-- 2. Drop new tables
DROP TABLE IF EXISTS issue_labels;
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS labels;
DROP TABLE IF EXISTS board_columns;
DROP TABLE IF EXISTS boards;

-- 3. Drop priority_type enum
DROP TYPE IF EXISTS priority_type;

-- 4. Recreate original issues table
CREATE TABLE IF NOT EXISTS issues (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id),
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      status_type NOT NULL DEFAULT 'Incomplete',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
