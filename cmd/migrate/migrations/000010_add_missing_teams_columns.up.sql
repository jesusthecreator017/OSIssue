ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS avatar_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS max_members INT NOT NULL DEFAULT 50;

-- Now that column exists (with no existing rows), make it NOT NULL
ALTER TABLE teams ALTER COLUMN created_by SET NOT NULL;
