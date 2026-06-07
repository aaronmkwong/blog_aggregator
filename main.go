package main

// _ blank identifier precedes github.com/lib/p
// because running this package's setup/initialization code, 
// but not going to use its code directly, so ignore unused import rule
import (
	"database/sql"
	"fmt"
	"os"
	"log"

	"github.com/aaronmkwong/blog_aggregator/internal/config"
	"github.com/aaronmkwong/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

// Holds shared application and database states
type state struct {
	db  *database.Queries
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

	// Open the connection
    db, err := sql.Open("postgres", cfg.DBURL) // using the dbURL from your config
    if err != nil {
        log.Fatalf("error connecting to database: %v", err)
    }
    defer db.Close() // It's good practice to close the db connection when main exits

    // Create SQLC queries instance
    dbQueries := database.New(db)

	// Initialize app state
	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Create command registry
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register login handler
	cmds.register("login", handlerLogin)

	// Register register handler
	cmds.register("register", handlerRegister)

	// Register register handler
	cmds.register("reset", handlerReset)

	// Register users handler
	cmds.register("users", handlerUsers)

	// Register agg handler
	cmds.register("agg", handlerAgg)	

	// Register add feed handler
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))

	// Register feed handler
	cmds.register("feeds", handlerFeed)

	// Register follow handler
	cmds.register("follow", middlewareLoggedIn(handlerFollow))

	// Register following handler
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

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
