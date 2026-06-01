-- This migration sets up our initial users table

-- WARNING: Do NOT modify or remove the "goose Up" and "goose Down" directives below!
-- Goose requires these exact comments to manage database migrations.

-- +goose Up
-- We use UUID for id to prevent ID enumeration attacks
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL, -- Tracks when the user registered
    updated_at TIMESTAMP NOT NULL, -- Tracks the last profile update
    name TEXT UNIQUE NOT NULL      -- Username must be unique and cannot be empty
);

-- +goose Down
-- This drops the users table to revert the migration
DROP TABLE users;