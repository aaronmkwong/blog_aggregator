// Defines the command and commands types, plus register and run methods

package main

import "fmt"

// Registers a command handler
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// Runs a command if it exists
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]

	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}
