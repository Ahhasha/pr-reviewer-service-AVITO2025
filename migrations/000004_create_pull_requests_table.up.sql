CREATE TABLE pull_requests (
    id TEXT PRIMARY KEY,                 
    name TEXT NOT NULL,                    
    author_id TEXT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    merged_at TIMESTAMPTZ,
    
    CONSTRAINT fk_pull_requests_author FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE INDEX idx_pull_requests_author ON pull_requests(author_id);
CREATE INDEX idx_pull_requests_status ON pull_requests(status);
CREATE INDEX idx_pull_requests_created_at ON pull_requests(created_at);