package main

import (
	"context"
	"fmt"
)

// handlerReset deletes all users and returns nil to indicate success
func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	// Execute sqlc deletion
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		// Return the error to let main.go print it and exit with code 1
		return fmt.Errorf("failed to delete users: %w", err)
	}

	// Report success to the user
	fmt.Println("Success: All users have been deleted from the database.")
	
	// Return nil to signal absolute success to the main execution loop
	return nil
}