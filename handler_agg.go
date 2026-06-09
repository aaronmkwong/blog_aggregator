package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"strings"

	"github.com/google/uuid"
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

		// Save posts from feed
		for _, item := range rssFeed.Channel.Item {

			// Initialize published_at as NULL
			publishedAt := sql.NullTime{}

			// Parse published time if present
			if item.PubDate != "" {

				parsedTime, err := parsePublishedTime(
					item.PubDate,
				)

				// Store parsed time if successful
				if err == nil {

					publishedAt = sql.NullTime{
						Time:  parsedTime,
						Valid: true,
					}

				} else {

					// Log parse error but continue saving post
					fmt.Printf(
						"Warning: unable to parse date '%s': %v\n",
						item.PubDate,
						err,
					)
				}
			}

			// Create post
			_, err = s.db.CreatePost(
				context.Background(),
				database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   now,
					UpdatedAt:   now,
					Title:       item.Title,
					Url:         item.Link,
					Description: sql.NullString{
						String: item.Description,
						Valid:  item.Description != "",
					},
					PublishedAt: publishedAt,
					FeedID:       feed.ID,
				},
			)

			// Ignore duplicate URLs
			if err != nil {

				// Skip duplicate posts
				if strings.Contains(
					err.Error(),
					"duplicate key value",
				) {
					continue
				}

				// Log other errors
				fmt.Printf(
					"Error creating post: %v\n",
					err,
				)
			}
		}
	}

	// Unreachable but required by compiler
	return nil
}

func parsePublishedTime(pubDate string) (time.Time, error) {
	// Try RFC1123Z first (e.g., "Mon, 02 Jan 2006 15:04:05 -0700")
	t, err := time.Parse(time.RFC1123Z, pubDate)
	if err == nil {
		return t, nil
	}

	// Fallback to RFC1123 if needed (e.g., "Mon, 02 Jan 2006 15:04:05 MST")
	return time.Parse(time.RFC1123, pubDate)
}