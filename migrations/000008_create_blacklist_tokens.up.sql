CREATE TABLE blacklist_tokens (
    id UUID PRIMARY KEY,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    expires TIMESTAMP NOT NULL,
    created TIMESTAMP NOT NULL,
    CONSTRAINT fk_blacklist_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_blacklist_tokens_token_hash ON blacklist_tokens(token_hash);
CREATE INDEX idx_blacklist_tokens_user_id ON blacklist_tokens(user_id);
CREATE INDEX idx_blacklist_tokens_expires_at ON blacklist_tokens(expires);
