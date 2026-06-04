package main

import (
    "context"
    "fmt"
)

// Handles the agg command
func handlerAgg(s *state, cmd command) error {

	// Fetch the test RSS feed
	feed, err := fetchFeed(
		context.Background(),
		"https://www.wagslane.dev/index.xml",
	)
	if err != nil {
		return err
	}

	// Print feed for debugging
	fmt.Printf("%+v\n", *feed)

	return nil
}