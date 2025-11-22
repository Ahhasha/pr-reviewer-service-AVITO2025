CREATE TABLE pr_reviewers (
    pr_id TEXT REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (pr_id, user_id),
    
    CONSTRAINT fk_pr_reviewers_pr FOREIGN KEY (pr_id) REFERENCES pull_requests(id),
    CONSTRAINT fk_pr_reviewers_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_pr_reviewers_user ON pr_reviewers(user_id);