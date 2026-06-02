package main

import (
	"context"
	"fmt"
	"os"
	"time"
	
	"github.com/google/uuid"
	"github.com/aaronmkwong/blog_aggregator/internal/database"
)

// Handles the login command
func handlerLogin(s *state, cmd command) error {

	// Ensure username was provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	// Get username from args
	username := cmd.args[0]

	// Verify user exists in database
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User %s does not exist\n", username)
		os.Exit(1)
	}

	// Save username to config
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	// Confirm user was set
	fmt.Printf("User set to %s\n", username)

	return nil
}

// Handles the register command
func handlerRegister(s *state, cmd command) error {

	// Ensure a username was provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	// Get username from args
	name := cmd.args[0]

	// Create the user in the database
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
		},
	)

	// Print error and exit if user creation fails
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}

	// Set the current user in config
	if err := s.cfg.SetUser(name); err != nil {
		return err
	}

	// Confirm user creation
	fmt.Printf("User %s was created\n", name)

	// Print user for debugging
	fmt.Println(user)

	return nil
}