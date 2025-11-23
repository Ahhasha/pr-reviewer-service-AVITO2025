CREATE TABLE teams (
    id TEXT PRIMARY KEY,
    team_name TEXT UNIQUE NOT NULL
);

CREATE INDEX idx_teams_name ON teams(team_name);