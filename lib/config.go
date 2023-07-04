package lib

type Config struct {
	// The bot token, used to log into discord.
	Token string `json:"discord"`

	// The command prefix used for commands.
	Prefix string `json:"command_prefix"`

	Categories CategoryConfig `json:"categories"`

	Roles RoleConfig `json:"roles"`
}

type CategoryConfig struct {
	// Channel categories that commands are enabled in. Commands sent in other categories are ignored.
	CommandsEnabled []string `json:"listening"`
}

type RoleConfig struct {
	// Roles for admin users. These users are allowed to send commands, and are treated specially in other areas (e.g. are never marked as spammers).
	AdminRoles []string `json:"admin"`
}
