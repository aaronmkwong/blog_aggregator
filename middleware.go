package main

import (
	"context"
	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Allows to change the function signature of handlers that require a logged in user to accept a user 
// as an argument and DRY up our code

// Loads the current user before running a handler
func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {

		// Get the current user from the database
		user, err := s.db.GetUser(
			context.Background(),
			s.cfg.CurrentUserName,
		)
		if err != nil {
			return err
		}

		// Run the wrapped handler
		return handler(s, cmd, user)
	}
}