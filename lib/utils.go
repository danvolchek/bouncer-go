package lib

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/danvolchek/bouncer-go/database"
	"github.com/hashicorp/golang-lru/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"regexp"
)

type Utils struct {
	Discord *discordgo.Session

	Config *Config

	Log zerolog.Logger

	DB *gorm.DB

	userIdToReplyThreadId *lruCache[string, string]
	replyThreadIdToUserId *lruCache[string, string]
}

// Reply sends a message in reply to another. It doesn't use the discord reply functionality.
func (u *Utils) Reply(replyTo *discordgo.Message, message string) {
	_, err := u.Discord.ChannelMessageSend(replyTo.ChannelID, message)
	if err != nil {
		u.Log.Error().Msgf("failed to send reply: %s", err)
	}
}

var snowflakeRegexp = regexp.MustCompile(`^\d+$`)

// UserFromId returns a user struct from a user id.
func (u *Utils) UserFromId(userId string) (*discordgo.User, error) {
	if !snowflakeRegexp.MatchString(userId) {
		return nil, errors.New("id isn't a snowflake")
	}

	user, errId := u.Discord.User(userId)
	return user, errId
}

// UserFromName returns a user struct from a user name. The user must be present in the guild provided, and the bot must have
// access to it's members.
func (u *Utils) UserFromName(userName, guildId string) (*discordgo.User, error) {
	guild, err := u.Discord.State.Guild(guildId)
	if err != nil {
		return nil, err
	}

	for _, member := range guild.Members {
		if member.User.Username == userName {
			return member.User, nil
		}
	}

	return nil, errors.New("no user with that name is in the server")
}

func (u *Utils) NewWithLog(funcs ...func(ctx zerolog.Context) zerolog.Context) *Utils {
	utils := *u

	ctx := utils.Log.With()
	for _, configure := range funcs {
		ctx = configure(ctx)
	}
	utils.Log = ctx.Logger()

	utils.DB = database.WithLogger(utils.DB, utils.Log)

	return &utils
}

func AddString(key, value string) func(c zerolog.Context) zerolog.Context {
	return func(c zerolog.Context) zerolog.Context { return c.Str(key, value) }
}

type lruCache[K comparable, V any] struct {
	*lru.Cache[K, V]
	valueFunc func(key K) (V, bool)
}

func (l *lruCache[K, V]) GetOrRetrieve(key K) (V, bool) {
	val, ok := l.Get(key)
	if ok {
		return val, true
	}

	val, ok = l.valueFunc(key)
	if ok {
		l.Add(key, val)
	}
	return val, ok
}
