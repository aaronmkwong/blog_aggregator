package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the addfeed command
func handlerAddFeed(s *state, cmd command) error {

	// Ensure name and URL were provided
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	// Get feed name from args
	name := cmd.args[0]

	// Get feed URL from args
	url := cmd.args[1]

	// Get the current user from the database
	user, err := s.db.GetUser(
		context.Background(),
		s.cfg.CurrentUserName,
	)
	if err != nil {
		return err
	}

	// Store current UTC timestamp
	now := time.Now().UTC()

	// Create the feed
	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return err
	}

	// Create a follow record for the feed creator
	_, err = s.db.CreateFeedFollow(
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

	// Confirm feed creation
	fmt.Printf("Feed '%s' created\n", feed.Name)

	return nil
}