package main

import (
	"context"
	"fmt"

	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the unfollow command
func handlerUnfollow(
	s *state,
	cmd command,
	user database.User,
) error {

	// Ensure a feed URL was provided
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	// Get feed URL from args
	url := cmd.args[0]

	// Look up the feed by URL
	feed, err := s.db.GetFeedByURL(
		context.Background(),
		url,
	)
	if err != nil {
		return err
	}

	// Remove the feed follow record
	err = s.db.UnfollowFeed(
		context.Background(),
		database.UnfollowFeedParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return err
	}

	// Confirm unfollow
	fmt.Printf(
		"%s unfollowed %s\n",
		user.Name,
		feed.Name,
	)

	return nil
}