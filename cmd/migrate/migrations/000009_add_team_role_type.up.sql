CREATE TYPE team_role AS ENUM ('owner', 'admin', 'member');

ALTER TABLE team_members
    ALTER COLUMN role DROP DEFAULT,
    ALTER COLUMN role SET DATA TYPE team_role USING role::team_role,
    ALTER COLUMN role SET DEFAULT 'member';
