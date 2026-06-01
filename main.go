package main

import (
	"fmt"
	"os"

	"github.com/aaronmkwong/blog_aggregator/internal/config"
)

// Holds shared application state
type state struct {
	cfg *config.Config
}

// Represents a CLI command
type command struct {
	name string
	args []string
}

// Stores registered command handlers
type commands struct {
	handlers map[string]func(*state, command) error
}

func main() {

	// Read config from disk
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize app state
	s := &state{
		cfg: &cfg,
	}

	// Create command registry
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register login handler
	cmds.register("login", handlerLogin)

	// Require command-line input
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	// Parse command name and args
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	// Execute command
	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
