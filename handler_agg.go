
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the agg command
func handlerAgg(s *state, cmd command) error {

	// Ensure a duration was provided
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	// Parse duration argument
	timeBetweenRequests, err := time.ParseDuration(
		cmd.args[0],
	)
	if err != nil {
		return err
	}

	// Print startup message
	fmt.Printf(
		"Collecting feeds every %s\n",
		timeBetweenRequests,
	)

	// Create ticker
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	// Run immediately and then on every tick
	for ; ; <-ticker.C {

		// Get the next feed to fetch
		feed, err := s.db.GetNextFeedToFetch(
			context.Background(),
		)
		if err != nil {
			fmt.Printf("Error getting next feed: %v\n", err)
			continue
		}

		// Store current UTC timestamp
		now := time.Now().UTC()

		// Mark feed as fetched
		err = s.db.MarkFeedFetched(
			context.Background(),
			database.MarkFeedFetchedParams{
				LastFetchedAt: sql.NullTime{
					Time:  now,
					Valid: true,
				},
				UpdatedAt: now,
				ID:        feed.ID,
			},
		)
		if err != nil {
			fmt.Printf("Error marking feed fetched: %v\n", err)
			continue
		}

		// Print feed being fetched
		fmt.Printf("Fetching feed: %s\n", feed.Name)

		// Fetch RSS feed
		rssFeed, err := fetchFeed(
			context.Background(),
			feed.Url,
		)
		if err != nil {
			fmt.Printf("Error fetching feed: %v\n", err)
			continue
		}

		// Print post titles
		for _, item := range rssFeed.Channel.Item {
			fmt.Printf(" - %s\n", item.Title)
		}
	}

	// Unreachable but required by compiler
	return nil
}
