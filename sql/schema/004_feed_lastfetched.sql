-- This migration sets up our initial users table

-- WARNING: Do NOT modify or remove the "goose Up" and "goose Down" directives below!
-- Goose requires these exact comments to manage database migrations.

-- This migration adds the last_fetched_at column to the feeds table to track when a feed was last scraped

-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;