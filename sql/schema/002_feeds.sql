-- This migration sets up our initial users table

-- WARNING: Do NOT modify or remove the "goose Up" and "goose Down" directives below!
-- Goose requires these exact comments to manage database migrations.

-- +goose Up
-- We use UUID for id to prevent ID enumeration attacks
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,                                     -- Tracks when the feed registered
    updated_at TIMESTAMP NOT NULL,                                     -- Tracks the last feed update
    name TEXT NOT NULL,                                                -- Name of the feed 
    url TEXT UNIQUE NOT NULL,                                          -- URL of the feed
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE       -- ID of the user who added this feed, cascade constraint to prevent orphaned records
);

-- +goose Down
-- This drops the users table to revert the migration
DROP TABLE feeds;