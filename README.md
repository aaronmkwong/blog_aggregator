**Use Case**

A multi-user CLI application needs to persist two pieces of state between runs: which database to connect to, and who is currently logged in. Rather than hardcoding these values or passing them as flags every time, the app reads and writes a JSON config file stored in the user's home directory.

Architecture
~/.gatorconfig.json          # persisted state (outside the project)

project root/
├── go.mod                   # module definition
├── main.go                  # entry point, orchestrates the app
└── internal/
    └── config/
        └── config.go        # config package: read/write JSON config

internal/ signals that the config package is private to this module — other external modules cannot import it.

main.go depends on the config package but the config package has no knowledge of main. Dependencies flow one way.

The JSON file on disk is the only persistence layer at this stage — no database yet.

**Coding Concepts**

Packages and visibility — exported identifiers (Config, Read, SetUser) are usable by main; unexported helpers (write, getConfigFilePath) are internal to the package.

Struct tags — json:"db_url" maps Go field names to JSON key names during marshalling and unmarshalling.

JSON encoding — json.Unmarshal deserializes JSON bytes into a Go struct; json.MarshalIndent serializes a struct back to formatted JSON bytes.

Methods vs functions — SetUser is a method on *Config (pointer receiver) so it can mutate the struct in place before writing to disk.

Error handling — Go's explicit error return pattern is used throughout, propagating errors up to main where they are logged and the program exits.

OS interaction — os.UserHomeDir, os.ReadFile, and os.WriteFile handle filesystem access in a cross-platform way.