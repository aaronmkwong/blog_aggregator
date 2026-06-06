package main

import (
	"context"
	"fmt"
)

// Handles the following command
func handlerFollowing(s *state, cmd command) error {

	// Get the current user
	user, err := s.db.GetUser(
		context.Background(),
		s.cfg.CurrentUserName,
	)
	if err != nil {
		return err
	}

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