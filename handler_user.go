package main

import "fmt"

// Handles the login command
func handlerLogin(s *state, cmd command) error {

	// Ensure username was provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	// Get username from args
	username := cmd.args[0]

	// Save username to config
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	// Confirm user was set
	fmt.Printf("User set to %s\n", username)

	return nil
}