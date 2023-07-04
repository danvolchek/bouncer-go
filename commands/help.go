package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/danvolchek/bouncer-go/lib"
	"github.com/danvolchek/bouncer-go/lib/components"
	"strconv"
	"strings"
)

const helpMessageTemplate = `
Issue a warning: {tick}{prefix}warn <user> <message>{tick}
Log a ban: {tick}{prefix}ban <user> <reason>{tick}
Log an unbanning: {tick}{prefix}unban <user> <reason>{tick}
Log a kick: {tick}{prefix}kick <user> <reason>{tick}
Ban with a pre-made scam message: {tick}{prefix}scam <user>{tick}
Preview what will be sent to the user {tick}{prefix}preview <warn/ban/kick> <reason>{tick}

Search for a user: {tick}{prefix}search <user>{tick}
Create a note about a user: {tick}{prefix}note <user> <message>{tick}
Remove a user's log: {tick}{prefix}remove <user> <index(optional)>{tick}
Edit a user's note: {tick}{prefix}edit <user> <index(optional)> <new_message>{tick}

Reply to a user in DMs: {tick}{prefix}reply <user> <message>{tick}
You can also Discord reply to a DM: {tick}{prefix}reply <message>{tick}
View users waiting for a reply: {tick}{prefix}waiting{tick}. Clear the list with {tick}{prefix}clear{tick}
Stop a user from sending DMs to us: {tick}{prefix}block/{prefix}unblock <user>{tick}

Sync bot commands to the server: {tick}{prefix}sync{tick}
Remove a user's 'Muted' role: {tick}{prefix}unmute <user>{tick}
Say a message as the bot: {tick}{prefix}say <channel> <message>{tick}

Watch a user's every move: {tick}{prefix}watch <user>{tick}
Remove user from watch list: {tick}{prefix}unwatch <user>{tick}
List watched users: {tick}{prefix}watchlist{tick}

Plot warn/ban stats: {tick}{prefix}graph{tick}
View bot uptime: {tick}{prefix}uptime{tick}

DMing users when they are banned is {tick}{dm_bans}{tick}
DMing users when they are warned is {tick}{dm_warns}{tick}
`

type help struct {
	replacer *strings.Replacer
}

func (h *help) Setup(utils *lib.Utils) {
	h.replacer = createReplacer(map[string]string{
		"tick":     "`",
		"prefix":   utils.Config.Prefix,
		"dm_bans":  strconv.FormatBool(utils.Config.DM.SendBanMessage),
		"dm_warns": strconv.FormatBool(utils.Config.DM.SendWarnMessage),
	})
}

func (h *help) Name() string {
	return "help"
}

func (h *help) Handle(command *components.CommandDetails, message *discordgo.Message, utils *lib.Utils) bool {
	utils.Reply(message, h.replacer.Replace(helpMessageTemplate))

	return true
}

func createReplacer(variables map[string]string) *strings.Replacer {
	i := 0
	args := make([]string, len(variables)*2)
	for key, value := range variables {
		args[i] = "{" + key + "}"
		args[i+1] = value
		i += 2
	}

	return strings.NewReplacer(args...)
}
