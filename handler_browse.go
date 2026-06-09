package main

import (
	"strconv"
	"fmt"
	"context"

	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the browse command
func handlerBrowse(
	s *state,
	cmd command,
	user database.User,
) error {

	// Default post limit
	limit := int32(2)

	// Parse optional limit
	if len(cmd.args) == 1 {

		n, err := strconv.Atoi(
			cmd.args[0],
		)
		if err != nil {
			return fmt.Errorf(
				"invalid limit",
			)
		}

		limit = int32(n)
	}

	// Get posts for current user
	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  limit,
		},
	)
	if err != nil {
		return err
	}

	// Print posts
	for _, post := range posts {

		fmt.Printf(
			"Title: %s\n",
			post.Title,
		)

		fmt.Printf(
			"URL: %s\n",
			post.Url,
		)

		if post.Description.Valid {
			fmt.Printf(
				"Description: %s\n",
				post.Description.String,
			)
		}

		fmt.Println()
	}

	return nil
}