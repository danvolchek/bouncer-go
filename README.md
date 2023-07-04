# Bouncer-go

It's https://github.com/aquova/bouncer written in go - mostly for fun.

WIP and everything is TBD

# Features complete so far

It's mostly framework stuff so far.

- Logging
- Load config file
- Add commands
- Interact with database

# Feature to be added

- Implementing all the commands
- All the background stuff that isn't a command
- Docker

# Running

Download the latest version of Go. On Arch: `sudo pacman -S go`.

The config file is the same as bouncer's, except it has an extra top level field `"database_path"` which should point to
your database file.

Then, one of:
 - rebuild every time: `go run main.go -config private/bouncer.db -debug`
 - build once, run multiple times: `go build -o bouncer main.go && ./bouncer -config private/bouncer.db -debug`
 - in IntelliJ IDE: run the `Run` run configuration. Config path is expected to be `private/bouncer.db`.

# Development

## Adding commands
- Add a file in the commands folder with a struct that implements the `Command` interface
- Add it to `commands.All`

## Adding new config fields
- Update `config.go`

## Adding new database fields
- Update `tables.go` as well as `database.allTables`

# Thoughts on the port

- discord-go is lower level, so more API calls are needed for tasks that discordpy would handle automatically
  - e.g. retrieving the roles a user has, or the category a channel is in
- everything being typed is already so much nicer than discord-py

# Thoughts on go and the libraries used
- zerolog is neat - has all the features I want out of a logger
  - previous logger framework experience I had was stdlib (lacking features) and logrus (verbose, hard to use)
- gorm looks neat so far - needs more usage to get it to work
