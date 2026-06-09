# Gator

Gator is a CLI RSS feed aggregator written in Go. It lets users register accounts, follow RSS feeds, continuously collect posts into PostgreSQL, and browse recent content from the feeds they follow.

## Requirements

Gator requires:

* Go
* PostgreSQL

Installation instructions:

* Go: https://go.dev/doc/install
* PostgreSQL: https://www.postgresql.org/download/

## Install

Install the CLI with:

```bash
go install github.com/aaronmkwong/blog_aggregator@latest
```

Replace the module path with your actual repository path.

## Configuration

Create a config file at:

```text
~/.gatorconfig.json
```

Example:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace the connection string with your PostgreSQL credentials.

The database schema should be created using the project's SQL migrations before running the application.

## Running Gator

Example workflow:

```bash
gator register alice
```

```bash
gator login alice
```

```bash
gator addfeed "Hacker News" https://news.ycombinator.com/rss
```

```bash
gator agg 1m
```

```bash
gator browse 10
```

## Useful Commands

### Users

```bash
gator register <username>
```

```bash
gator login <username>
```

```bash
gator users
```

### Feeds

```bash
gator addfeed "<feed_name>" <feed_url>
```

```bash
gator feeds
```

```bash
gator follow <feed_url>
```

```bash
gator unfollow <feed_url>
```

```bash
gator following
```

### Aggregation

```bash
gator agg <duration>
```

Example:

```bash
gator agg 30s
```

### Browse Posts

```bash
gator browse
```

```bash
gator browse <limit>
```

## Architecture

~/.gatorconfig.json — persisted state (outside the project)

```
project root/
├── go.mod                             # module definition
├── go.sum                             # dependency checksums
├── main.go                            # entry point: reads config, initializes state, registers/runs commands
├── middleware.go                      # enables handlers function signatures requiring logged in user to accept user as an argument 
├── commands.go                        # command system: command structs, registry, run/register methods
├── handler_user.go                    # user-related command handlers: login, register
├── handler_users.go                   # users command handler: lists users and marks the current user
├── handler_reset.go                   # reset command handler: deletes all users (dev/testing utility)
├── handler_agg.go                     # agg command handler: fetches and prints RSS feed
├── handler_browse.go                  # ...
├── handler_addfeed.go                 # addfeed command handler: creates a feed and auto-follows it
├── handler_feed.go                    # feed command handler: prints all feeds in the database
├── handler_follow.go                  # follow command handler: creates a feed follow for the current user
├── handler_following.go               # following command handler: prints all feeds the current user follows
├── handler_unfollow.go                # unfollow command handler: unfollows for the current user
├── rss_feed.go                        # RSS types (RSSFeed, RSSItem) and fetchFeed function
├── sqlc.yaml                          # SQLC config: maps schema/queries dirs to generated Go output
├── README.md                          # project documentation
├── sql/
│   ├── schema/
│   │   ├── 001_users.sql              # Goose migration: creates/drops the users table (up/down)
│   │   ├── 002_feeds.sql              # Goose migration: creates/drops the feeds table (up/down) with ON DELETE CASCADE
│   │   └── 003_feed_follows.sql       # Goose migration: creates/drops the feed_follows table (up/down) with ON DELETE CASCADE
│   │   └── 004_feed_lastfetched.sql   # Goose migration: adds the last_fetched_at column to the feeds table to track when a feed was last scraped
│   │   └── 005_posts.sql              # Goose migration: creates/drops the posts table (up/down)
│   └── queries/
│       ├── users.sql                  # SQLC query: CreateUser, GetUser
│       ├── get_users.sql              # SQLC query: GetUsers
│       ├── del_users.sql              # SQLC query: DeleteUsers
│       ├── feeds.sql                  # SQLC query: CreateFeed, GetFeeds, GetFeedByURL
│       ├── feed_follow.sql            # SQLC query: CreateFeedFollow
│       └── get_feed_follows_user.sql  # SQLC query: GetFeedFollowsForUser
│       └── feed_unfollow.sql          # SQLC query: UnfollowFeed
│       └── feed_mark_fetched.sql      # SQLC query: MarkFeedFetched
│       └── feed_next_fetch.sql        # SQLC query: GetNextFeedToFetch
│       └── get_posts.sql              # SQLC query: GetPostsForUser
│       └── posts.sql                  # SQLC query: CreatePost
└── internal/
    ├── config/
    │   └── config.go                  # config package: read/write JSON config and update current user
    └── database/                      # generated by SQLC — do not edit manually
        ├── db.go                      # database connection wrapper (Queries struct)
        ├── models.go                  # Go structs mirroring database table rows (User, Feed, FeedFollow)
        ├── users.sql.go               # generated from users.sql (CreateUser, GetUser)
        ├── get_users.sql.go           # generated from get_users.sql (GetUsers)
        ├── del_users.sql.go           # generated from del_users.sql (DeleteUsers)
        ├── feeds.sql.go               # generated from feeds.sql (CreateFeed, GetFeeds, GetFeedByURL)
        ├── feed_follow.sql.go         # generated from feed_follow.sql (CreateFeedFollow)
        └── get_feed_follows_user.sql.go  # generated from get_feed_follows_user.sql (GetFeedFollowsForUser)
        └── feed_unfollow.sql.go       # generated from feed_unfollow.sql (UnfolowFeed)
        └── feed_markfetched.sql.go    # generated from feed_markfetched.sql (MarkFeedFetched)
        └── feed_next_fetch.sql.go     # generated from feed_next_fetch.sql (GetNextFeedToFetch)
        └── get_posts.go               # generated from get_posts.sql (GetPostsForUser)
        └── posts.sql.go               # generated from posts.sql (CreatePost)
```

internal/ signals that both config and database packages are private to this module.

main.go depends on both config and database packages; neither has knowledge of main. Dependencies flow one way.

The PostgreSQL database is the persistence layer for users. The JSON file persists the currently logged-in username.

main.go                          -> startup: read config, open DB connection, create state,
                                    register commands (login, register, reset, users, agg,
                                    addfeed, follow, following), parse os.Args

middleware.go                    -> middleware: wrap handlers that require a logged-in user
                                    by fetching the current user from the DB before
                                    delegating to the handler

commands.go                      -> command, commands, run/register

handler_user.go                  -> handlerLogin, handlerRegister

handler_users.go                 -> handlerUsers: calls GetUsers, prints each username,
                                    and marks the currently logged-in user

handler_reset.go                 -> handlerReset: calls DeleteUsers query, dev/test utility only

handler_agg.go                   -> handlerAgg: calls fetchFeed with a hardcoded URL,
                                    prints the resulting RSSFeed struct to the console

handler_browse.go                -> handlerBrowse: retrieves the posts from feeds followed by the 
                                    logged-in user up to an optional limit, and prints their details 
                                    (title, URL, description) to the console

handler_addfeed.go               -> handlerAddFeed: checks arguments, gets the current user,
                                    creates a new feed, auto-creates a feed follow record
                                    for the creator, and prints confirmation

handler_feed.go                  -> handlerFeed: retrieves all feeds from the database,
                                    looks up the user associated with each feed, and prints
                                    each feed's name, URL, and creator

handler_follow.go                -> handlerFollow: checks arguments, gets the current user,
                                    looks up the feed by URL, creates a feed follow record,
                                    and prints the user and feed name

handler_unfollow.go              -> handlerUnfollow: checks arguments, gets the current user, 
                                    looks up the feed by URL, deletes the matching feed follow record 
                                    for that user and feed, and prints a confirmation message                                    

handler_following.go             -> handlerFollowing: gets the current user, retrieves all
                                    feed follows for that user, and prints each feed name

rss_feed.go                      -> RSSFeed, RSSItem structs with xml tags
                                    fetchFeed: creates HTTP request with context and User-Agent header,
                                    executes request, reads body, unmarshals XML, unescapes HTML entities
                                    in channel and item Title/Description fields

sql/schema/                      -> Goose migrations (schema management)
                                    001_users: manages users table
                                    002_feeds: manages feeds table with foreign key referencing users(id)
                                    003_feed_follows: manages feed_follows table with foreign keys
                                    referencing users(id) and feeds(id), both ON DELETE CASCADE,
                                    and a unique constraint on (user_id, feed_id) pairs

sql/queries/                     -> SQLC query definitions (raw SQL)
                                    :one  — single row (CreateUser, GetUser, CreateFeed, GetFeedByURL, CreateFeedFollow)
                                    :many — multiple rows (GetUsers, GetFeeds, GetFeedFollowsForUser)
                                    :exec — no returned rows (DeleteUsers)

internal/database                -> SQLC-generated type-safe Go code (not for editing): provides typed
                                    interfaces to interact with database records securely

## Architecture<br>

\*Add sorting and filtering options to the browse command<br>
\*Add pagination to the browse command<br>
\*Add concurrency to the agg command so that it can fetch more frequently<br>
\*Add a search command that allows for fuzzy searching of posts<br>
\*Add bookmarking or liking posts<br>
\*Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)<br>
\*Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely<br>
\*Write a service manager that keeps the agg command running in the background and restarts it if it crashes<br>