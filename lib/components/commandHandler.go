package components

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/danvolchek/bouncer-go/database"
	"github.com/danvolchek/bouncer-go/lib"
	uuid2 "github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
	"strings"
)

// CommandHandler is a component responsible for handling commands for the bot - i.e. messages that begin with the command prefix.
// The handler is a generic orchestrator capable of executing any command - it handles common stuff like making sure a command
// is being sent in the first place, the user has the authority to send the command, parsing the command text, etc.
// Commands themselves should implement the Command interface and be passed when creating the handler.
type CommandHandler struct {
	// map from command name to command, for easy access
	commands map[string]Command

	// general utilities
	*lib.Utils
}

// Command is the interface commands should implement.
type Command interface {
	// Name should return the name of the command, case-sensitive.
	Name() string

	// Handle should perform the action the command does.
	// It should return a boolean indicating whether the command succeeded - if not the handler will reply with a
	// UUID to look up the logs. The command should log anything relevant, including errors that happen.
	Handle(command *CommandDetails, message *discordgo.Message, utils *lib.Utils) bool
}

// CommandDetails provides some pre-processed information about a command that was executed for convenience.
// Full details can be found in the message parameter of Command::Handle.
type CommandDetails struct {
	// Name is the name of the command.
	Name string

	// Args is the command arguments.
	Args []string
}

func (c CommandDetails) ShortString() string {
	if len(c.Args) == 0 {
		return c.Name
	}

	str := c.Name + " " + strings.Join(c.Args, " ")
	if len(str) > 30 {
		str = str[0:30]
	}

	return str
}

// NewCommandHandler creates a command handler that runs the provided commands.
func NewCommandHandler(commands []Command) (*CommandHandler, error) {
	var commandMap = make(map[string]Command, len(commands))

	for _, command := range commands {
		name := command.Name()
		if _, ok := commandMap[name]; ok {
			return nil, fmt.Errorf("duplicate command '%s'", name)
		}

		commandMap[name] = command
	}

	return &CommandHandler{commands: commandMap}, nil
}

// Setup is called before the bot is started. Register any hooks/perform any initial setup here.
func (c *CommandHandler) Setup(utils *lib.Utils) {
	c.Utils = utils

	c.Discord.AddHandler(c.handleCommand)
}

// handleCommand is called on every new message being sent and runs the appropriate command based on the message text.
func (c *CommandHandler) handleCommand(_ *discordgo.Session, messageCreate *discordgo.MessageCreate) {
	// Add uuid to logs about this command execution for easier debugging
	uuid := uuid2.New().String()
	log := c.Log.With().Str("uuid", uuid).Logger()

	message := messageCreate.Message

	// Check if this message should be ignored or not - is the user an admin, in the right channel, etc
	if c.shouldIgnoreMessage(message, log) {
		return
	}

	// From here on, provide more error details to the user because they're an admin
	sendUUID := func() {
		c.Reply(message, fmt.Sprintf("Oops, something went wrong handling that message. Check the logs for `%s`.", uuid))
	}

	commandDetails, err := c.parseCommand(message)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse message")
		sendUUID()
		return
	}

	command, ok := c.commands[commandDetails.Name]
	if !ok {
		c.Reply(message, fmt.Sprintf("Unknown command `%s` - see `%shelp`", commandDetails.Name, c.Config.Prefix))
		return
	}

	// Now that the command is parsed, add the command and args to logs
	log = log.With().
		Str("command", commandDetails.ShortString()).
		Logger()

	utils := &(*c.Utils)
	utils.Log = log
	utils.DB = database.WithLogger(utils.DB, log)

	defer func() {
		if r := recover(); r != nil {
			sendUUID()
		}
	}()

	ok = command.Handle(commandDetails, message, utils)
	if !ok {
		log.Warn().Msg("command failed")
		sendUUID()
	}
}

// shouldIgnoreMessage returns whether the message that was sent should be ignored
// any errors should be logged but not sent to the user who sent the message because that user may not be an admin
func (c *CommandHandler) shouldIgnoreMessage(message *discordgo.Message, log zerolog.Logger) bool {
	// Ignore messages sent by this user or other bots
	if message.Author.ID == c.Discord.State.User.ID || message.Author.Bot {
		return true
	}

	// Ignore messages in non-enabled channels
	{
		channel, err := c.Discord.Channel(message.ChannelID)
		if err != nil {
			log.Error().Err(err).Msg("ignoring message - failed to retrieve channel info")
			return true
		}

		if !slices.Contains(c.Config.Categories.CommandsEnabled, channel.ParentID) {
			log.Trace().
				Str("category", channel.ParentID).
				Str("channel", message.ChannelID).
				Msg("ignoring message in non-enabled category")
			return true
		}
	}

	// Ignore messages from non-admin users
	{
		user, err := c.Discord.State.Member(message.GuildID, message.Author.ID)
		if err != nil {
			log.Error().Err(err).Msg("ignoring message - failed to retrieve member info")
			return true
		}

		if !slices.ContainsFunc(c.Config.Roles.AdminRoles, func(adminRole string) bool {
			return slices.Contains(user.Roles, adminRole)
		}) {
			log.Debug().Str("user", message.Author.Username).Msg("ignoring message from non-admin user")
			return true
		}

	}

	// Ignore messages without bot prefix
	if !strings.HasPrefix(strings.TrimSpace(message.Content), c.Config.Prefix) {
		log.Debug().Msg("ignoring message without bot prefix")
		return true
	}

	return false
}

func (c *CommandHandler) parseCommand(message *discordgo.Message) (*CommandDetails, error) {
	// Strip prefix from message (it must be present because of the above check to ignore messages without it)
	messageContent := strings.TrimSpace(message.Content)[len(c.Config.Prefix):]

	// Parse command name and args out of the message
	commandName, rawArgs, ok := strings.Cut(messageContent, " ")
	if !ok {
		// If there are no spaces, the command has no args
		return &CommandDetails{
			Name: messageContent,
			Args: nil,
		}, nil
	}

	// Trim spaces in commands to allow for extra spaces
	commandName = strings.TrimSpace(commandName)

	// Split args into components and also trim them
	commandArgs := strings.Split(rawArgs, " ")
	for i, arg := range commandArgs {
		commandArgs[i] = strings.TrimSpace(arg)
	}

	return &CommandDetails{
		Name: commandName,
		Args: commandArgs,
	}, nil
}
