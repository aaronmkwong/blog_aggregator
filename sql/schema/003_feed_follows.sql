-- WARNING: Do NOT modify or remove the "goose Up" and "goose Down" directives below!
-- Goose requires these exact comments to manage database migrations.

-- +goose Up
-- We use UUID for id to prevent ID enumeration attacks
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
-- This drops the users table to revert the migration
DROP TABLE feed_follows;