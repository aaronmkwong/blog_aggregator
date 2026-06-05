package main

import (
	"context"
	"fmt"
)

// Handles the feeds command
func handlerFeed(s *state, cmd command) error {

	// Get all feeds from database
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Print each feed
	for _, feed := range feeds {
		fmt.Printf(
			"Name: %s\nURL: %s\nUser: %s\n\n",
			feed.Name,
			feed.Url,
			feed.UserName,
		)
	}

	return nil
}