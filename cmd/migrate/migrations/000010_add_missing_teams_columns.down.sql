ALTER TABLE teams
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS max_members;
