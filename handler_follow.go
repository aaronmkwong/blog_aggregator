package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the follow command
func handlerFollow(s *state, cmd command) error {

	// Ensure a feed URL was provided
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	// Get feed URL from args
	url := cmd.args[0]

	// Get the current user
	user, err := s.db.GetUser(
		context.Background(),
		s.cfg.CurrentUserName,
	)
	if err != nil {
		return err
	}

	// Look up the feed by URL
	feed, err := s.db.GetFeedByURL(
		context.Background(),
		url,
	)
	if err != nil {
		return err
	}

	// Store current UTC timestamp
	now := time.Now().UTC()

	// Create the feed follow record
	follow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}

	// Print the created follow relationship
	fmt.Printf(
		"%s is now following %s\n",
		follow.UserName,
		follow.FeedName,
	)

	return nil
}