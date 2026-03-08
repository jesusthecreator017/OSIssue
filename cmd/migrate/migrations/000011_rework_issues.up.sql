-- 000011_rework_issues.up.sql
-- Rework issues to support boards, columns, priorities, labels, assignees

-- 1. Create boards table
CREATE TABLE IF NOT EXISTS boards (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    owner_user_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    owner_team_id   UUID REFERENCES teams(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT boards_owner_check CHECK (
        (owner_user_id IS NOT NULL AND owner_team_id IS NULL) OR
        (owner_user_id IS NULL AND owner_team_id IS NOT NULL)
    )
);

-- 2. Create board_columns table
CREATE TABLE IF NOT EXISTS board_columns (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id    UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    position    INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(board_id, position)
);

-- 3. Create priority_type enum
CREATE TYPE priority_type AS ENUM ('Low', 'Medium', 'High', 'Critical');

-- 4. Create labels table
CREATE TABLE IF NOT EXISTS labels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    color       TEXT NOT NULL DEFAULT '#6b7280',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 5. Drop old issues table (BIGSERIAL version)
DROP TABLE IF EXISTS issues;

-- 6. Recreate issues table with UUID and new fields
CREATE TABLE issues (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    assignee_id     UUID REFERENCES users(id),
    team_id         UUID REFERENCES teams(id) ON DELETE SET NULL,
    board_column_id UUID REFERENCES board_columns(id) ON DELETE SET NULL,
    position        INT NOT NULL DEFAULT 0,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    priority        priority_type NOT NULL DEFAULT 'Medium',
    due_date        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 7. Create issue_labels join table
CREATE TABLE IF NOT EXISTS issue_labels (
    issue_id    UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    label_id    UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, label_id)
);

-- 8. Seed default labels
INSERT INTO labels (name, color) VALUES
    ('bug', '#ef4444'),
    ('feature', '#3b82f6'),
    ('enhancement', '#8b5cf6'),
    ('documentation', '#6b7280'),
    ('urgent', '#f97316')
ON CONFLICT (name) DO NOTHING;

-- 9. Drop the old status_type enum
DROP TYPE IF EXISTS status_type;
