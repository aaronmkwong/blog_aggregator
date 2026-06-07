package main

import (
	"context"
	"fmt"

	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the following command
func handlerFollowing(s *state, cmd command, user database.User) error {

	// Get all feed follows for the current user
	follows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		return err
	}

	// Print each followed feed name
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}

	return nil
}