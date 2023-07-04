package lib

type Config struct {
	// The bot token, used to log into discord.
	Token string `json:"discord"`

	// The command prefix used for commands.
	Prefix string `json:"command_prefix"`

	Servers ServerConfig `json:"servers"`

	Categories CategoryConfig `json:"categories"`

	Channels ChannelConfig `json:"channels"`

	Roles RoleConfig `json:"roles"`

	Users UserConfig `json:"users"`

	DM DMConfig `json:"DM"`
}

type ServerConfig struct {
	// Server the bot is for. The bot can only run for one server.
	Home string `json:"home"`
}

type CategoryConfig struct {
	// Categories where commands are enabled in. Commands sent in other categories are ignored.
	CommandsEnabled []string `json:"listening"`
}

type ChannelConfig struct {
	// Channel to log user DMs to.
	Mailbox string `json:"mailbox"`

	// Channel to log spam notifications to.
	Spam string `json:"spam"`

	// Channels where spam is ignored/allowed.
	SpamIgnored []string `json:"ignore_spam"`

	// Channel to log general notifications to.
	Log string `json:"log"`

	// Channel to log system notifications (?) to.
	SysLog string `json:"syslog"`

	// Channel to log notifications about watched users to.
	Watchlist string `json:"watchlist"`

	// Channel to log ban appeal notifications to.
	BanAppeal string `json:"ban_appeal"`
}

type RoleConfig struct {
	// Roles for admin users. These users are allowed to send commands, and are treated specially in other areas (e.g. are never marked as spammers).
	Admin []string `json:"admin"`

	// Roles for bot owners. These users are also allowed to use the bot in certain ways.
	Owner []string `json:"owner"`

	// Roles to add to DM threads.
	DMThread []string `json:"dm_threads"`
}

type UserConfig struct {
	// Bot owners, able to use the bot/change certain settings.
	Owners []string `json:"owners"`
}

type DMConfig struct {
	// Whether to send users a message when they're banned indicating why.
	SendBanMessage bool `json:"ban"`

	// Whether to send users a message when they're warned indicating why.
	SendWarnMessage bool `json:"warn"`
}
