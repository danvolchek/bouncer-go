package components

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/danvolchek/bouncer-go/lib"
	uuid2 "github.com/google/uuid"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

//go:generate mockgen -typed -destination mocks/mock_command.go . Command

// Commands is a component responsible for handling commands for the bot - i.e. messages that begin with the command prefix.
// The handler is a generic orchestrator capable of executing any command - it handles common stuff like making sure a command
// is being sent in the first place, the user has the authority to send the command, parsing the command text, etc.
// Commands themselves should implement the Command interface and be passed when creating the handler.
type Commands struct {
	// map from command name to command, for easy access
	commands map[string]Command

	// general utilities
	*lib.Utils
}

// Command is the interface commands should implement.
type Command interface {
	// Name should return the name of the command, case-sensitive.
	Name() string

	// RequiresUser should return whether the command requires a user argument. If true and not provided, the handler
	// will reply with an error.
	RequiresUser() bool

	// Setup is called before the bot is started, after configs/the db is loaded. Perform initial setup here.
	// Don't save the utils class.
	Setup(utils *lib.Utils)

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

	// User is first user reference in the command. Will be nil if there are none.
	User *discordgo.User
}

func (c CommandDetails) ShortString() string {
	if len(c.Args) == 0 {
		return c.Name
	}

	str := c.Name + " " + strings.Join(c.Args, " ")
	if len(str) > 30 {
		str = str[0:27] + "..."
	}

	return str
}

// NewCommands creates a command handler that runs the provided commands.
func NewCommands(commands []Command) (*Commands, error) {
	var commandMap = make(map[string]Command, len(commands))

	for _, command := range commands {
		name := command.Name()
		if _, ok := commandMap[name]; ok {
			return nil, fmt.Errorf("duplicate command '%s'", name)
		}

		commandMap[name] = command
	}

	return &Commands{commands: commandMap}, nil
}

// Setup is called before the bot is started. Register any hooks/perform any initial setup here.
func (c *Commands) Setup(utils *lib.Utils) {
	c.Utils = utils

	for name, command := range c.commands {
		c.Log.Debug().Str("command", name).Msg("Registering command")
		command.Setup(c.NewWithLog(lib.AddString("command", name)))
	}

	c.Discord.AddHandler(c.handleCommand)
}

// handleCommand is called on every new message being sent and runs the appropriate command based on the message text.
func (c *Commands) handleCommand(_ *discordgo.Session, messageCreate *discordgo.MessageCreate) {
	// Add uuid to logs about this command execution for easier debugging
	uuid := uuid2.New().String()

	invoker := commandInvoker{
		uuid:     uuid,
		message:  messageCreate.Message,
		commands: c.commands,
		Utils:    c.Utils.NewWithLog(lib.AddString("uuid", uuid)),
	}

	invoker.invoke()
}

// commandInvoker invokes a command once.
type commandInvoker struct {
	// uuid to trace execution of this command
	uuid string

	// message including the command name + args
	message *discordgo.Message

	// commands to lookup implementation from
	commands map[string]Command

	// general utilities
	*lib.Utils
}

// invoke invokes the command
func (c commandInvoker) invoke() {
	if c.shouldIgnoreMessage() {
		return
	}

	if c.Log.Debug().Enabled() {
		c.sendUUID(false)
	}

	commandDetails, err := c.parseCommand()
	if err != nil {
		c.Log.Error().Err(err).Msg("failed to parse message")
		c.sendUUID(true)
		return
	}

	command, ok := c.commands[commandDetails.Name]
	if !ok {
		c.Log.Debug().Str("name", commandDetails.Name).Msg("no command exists with this name")
		c.Reply(c.message, fmt.Sprintf("Unknown command `%s` - see `%shelp`", commandDetails.Name, c.Config.Prefix))
		return
	}

	c.Utils = c.NewWithLog(lib.AddString("command", commandDetails.ShortString()))

	if command.RequiresUser() {
		ok := c.setUser(commandDetails)
		if !ok {
			c.Log.Warn().Msg("user not found, but one is required")
			c.Reply(c.message, fmt.Sprintf("This command requires a valid user - see `%shelp`", c.Config.Prefix))
			return
		}
	}

	defer func() {
		if r := recover(); r != nil {
			c.Log.Error().Any("panic", r).Msg("command panicked")
			c.sendUUID(true)
		}
	}()

	ok = command.Handle(commandDetails, c.message, c.Utils)
	if !ok {
		c.Log.Warn().Msg("command failed")
		c.sendUUID(true)
	}
}

// sendUUID sends a discord message with the invoker's uuid
func (c commandInvoker) sendUUID(wasError bool) {
	if !wasError {
		c.Reply(c.message, fmt.Sprintf("UUID for logs is `%s`.", c.uuid))
	} else {
		c.Reply(c.message, fmt.Sprintf("Oops, something went wrong handling that message. Check the logs for `%s`.", c.uuid))
	}
}

// shouldIgnoreMessage returns whether the message that was sent should be ignored
// any errors should be logged but not sent to the user who sent the message because that user may not be an admin
func (c commandInvoker) shouldIgnoreMessage() bool {
	// Ignore messages sent by this user or other bots
	if c.message.Author.ID == c.Discord.State.User.ID || c.message.Author.Bot {
		return true
	}

	// Ignore messages in non-enabled channels
	{
		channel, err := c.Discord.State.Channel(c.message.ChannelID)
		if err != nil {
			c.Log.Error().Err(err).Msg("ignoring message - failed to retrieve channel info")
			return true
		}

		if !slices.Contains(c.Config.Categories.CommandsEnabled, channel.ParentID) {
			c.Log.Trace().
				Str("category", channel.ParentID).
				Str("channel", c.message.ChannelID).
				Msg("ignoring message in non-enabled category")
			return true
		}
	}

	// Ignore messages from non-admin users
	{
		user, err := c.Discord.State.Member(c.message.GuildID, c.message.Author.ID)
		if err != nil {
			c.Log.Error().Err(err).Msg("ignoring message - failed to retrieve member info")
			return true
		}

		if !slices.ContainsFunc(c.Config.Roles.Admin, func(adminRole string) bool {
			return slices.Contains(user.Roles, adminRole)
		}) {
			c.Log.Debug().Str("user", c.message.Author.Username).Msg("ignoring message from non-admin user")
			return true
		}

	}

	// Ignore messages without bot prefix
	if !strings.HasPrefix(strings.TrimSpace(c.message.Content), c.Config.Prefix) {
		c.Log.Debug().Msg("ignoring message without bot prefix")
		return true
	}

	return false
}

// parseCommand parses the raw message text into a command
func (c commandInvoker) parseCommand() (*CommandDetails, error) {
	// Strip prefix from message (it must be present because of the above check to ignore messages without it)
	messageContent := strings.TrimSpace(c.message.Content)[len(c.Config.Prefix):]

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

// setUser updates command details with a user, returning whether it was successful
func (c commandInvoker) setUser(command *CommandDetails) bool {
	// First, parse user from command arguments
	user := c.getUserFromArg(command, c.message.GuildID)

	if user != nil {
		// If found, update args and user
		command.Args = command.Args[1:]
		command.User = user

		return true
	}

	// TODO: otherwise, try checking if it's a reply thread
	return false
}

var userPingRegexp = regexp.MustCompile(`<@(\d+)>`)

// getUserFromArg returns a user from command arguments
func (c commandInvoker) getUserFromArg(command *CommandDetails, guildId string) *discordgo.User {
	if len(command.Args) == 0 {
		c.Log.Debug().
			Msg("command has no args to get user mention from")

		return nil
	}

	userRef := command.Args[0]

	if match := userPingRegexp.FindStringSubmatch(userRef); match != nil {
		userRef = match[1]
	}

	// first try getting the user from their id directly, or a mention
	user, errId := c.UserFromId(userRef)
	if errId == nil {
		return user
	}

	// next try getting the user from a name
	user, errName := c.UserFromName(userRef, guildId)
	if errName == nil {
		return user
	}

	c.Log.Debug().Str("ref", userRef).Err(errId).Msg("cmd arg isn't a valid user id")

	c.Log.Debug().Str("ref", userRef).Err(errName).Msg("cmd arg isn't a valid user name")

	// otherwise there's no user
	return nil
}
