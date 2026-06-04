package main

import (
	"context"
	"fmt"
)

// Handles the users command
func handlerUsers(s *state, cmd command) error {

	// Get all usernames from database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	// Print each username
	for _, user := range users {

		// Mark the currently logged-in user
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
			continue
		}

		// Print non-current users
		fmt.Printf("* %s\n", user)
	}

	return nil
}