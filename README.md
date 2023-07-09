# Bouncer-go

It's https://github.com/aquova/bouncer written in go - mostly for fun.

WIP and everything is TBD

# Features complete so far

It's mostly framework stuff so far.

- Logging
- Load config file
- Add commands
- Interact with database

## Commands

| Command     | Implemented |
|-------------|-------------|
| `ban`       | ❌           |
| `block`     | ❌           |
| `clear`     | ❌           |
| `edit`      | ❌           |
| `graph`     | ❌           |
| `help`      | ✅           |
| `kick`      | ❌           |
| `note`      | ❌           |
| `preview`   | ❌           |
| `remove`    | ❌           |
| `reply`     | ❌           |
| `say`       | ❌           |
| `scam`      | ❌           |
| `search`    | ❌           |
| `sync`      | ❌           |
| `unban`     | ❌           |
| `unblock`   | ❌           |
| `uptime`    | ❌           |
| `waiting`   | ❌           |
| `warn`      | ❌           |
| `watch`     | ❌           |
| `watchlist` | ❌           |
| `unmute`    | ❌           |
| `unwatch`   | ❌           |


# Feature to be added

- Implementing all the commands
- All the background stuff that isn't a command
- Docker

# Running

Download the latest version of Go. On Arch: `sudo pacman -S go`.

Also get [task](https://taskfile.dev/), an alternative to make. On Arch: `yay -S go-task-bin`. I use `alias task="go-task"`.

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
 - CLI: `task run`
 - In IntelliJ IDE: run the `Run` run configuration. Configs are expected to be in `private/`.

To just build the binary, run `task build`.

If you want to run the go commands directly, see [taskfile.yml](taskfile.yml).

# Development

## Setting up unit tests
Run `task generate` to generate the mocks. These aren't committed.

## Adding commands
- Add a file in the commands folder with a struct that implements the `Command` interface
- Add it to `commands.All`

## Adding new config fields
- Update `config.go`

## Adding new database fields
- Update `tables.go` as well as `database.allTables`

## Docs
- Discordgo: https://github.com/bwmarrin/discordgo, https://pkg.go.dev/github.com/bwmarrin/discordgo
- Zerolog: https://github.com/rs/zerolog
- Gomock: https://github.com/uber/mock
- Other Go libraries: https://github.com/eleven26/awesome-go-stars
- Slices: https://pkg.go.dev/golang.org/x/exp/slices
- Discordpy: https://discordpy.readthedocs.io/en/latest/api.html
- Taskfile: https://taskfile.dev/

# Thoughts on the port
- discordgo is lower level, so more API calls are needed for tasks that discordpy would handle automatically
  - e.g. retrieving the roles a user has, or the category a channel is in
- everything being typed is already so much nicer than discord-py

# Thoughts on go and the libraries used
- zerolog is neat - has all the features I want out of a logger
  - previous logger framework experience I had was stdlib (lacking features) and logrus (verbose, hard to use)
- gorm looks neat so far - needs more usage to get it to work
