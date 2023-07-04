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

The program takes in a config directory and expects there to be a `bouncer.db` and a `config.json` file inside of them.

The config file is similar to bouncer's except:
 - All ids are strings instead of numbers because discord-go uses strings for ids rather than numbers like discordpy
 - The `owner` user id field is moved into the `roles` config named `owner` (list of role ids)
 - The `DM` configs are booleans rather than numbers
 - The `rolesToAddToThreads` field was moved under `roles` named `dm_threads` and it's parent `messageForwarding` removed
 - The `debug` field isn't used

Sample `config.json` file (see [config.go](lib/config.go) for meaning):
```json
{
  "discord": "xxxxxx",
  "command_prefix": "$",
  "servers": {
    "home": "12345678910"
  },
  "categories": {
    "listening": [
      "12345678910"
    ]
  },
  "channels": {
    "mailbox": "12345678910",
    "spam": "12345678910",
    "ignore_spam": [],
    "log": "12345678910",
    "syslog": "12345678910",
    "watchlist": "12345678910",
    "ban_appeal": "12345678910"
  },
  "roles": {
    "admin": [
      "12345678910"
    ],
    "owner": [
      "12345678910"
    ],
    "dm_thread": [
      "12345678910"
    ]
  },
  "DM":{
    "ban": true,
    "warn": true
  }
}
```

The database file is identical to bouncer's.

Then, one of:
 - Rebuild every time: `go run main.go -config private/ -debug`
 - Build once, run multiple times: `go build -o bouncer main.go && ./bouncer -config private/ -debug`
 - In IntelliJ IDE: run the `Run` run configuration. Configs are expected to be in `private/`.

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
