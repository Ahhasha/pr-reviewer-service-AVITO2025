CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username NEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);