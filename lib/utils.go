package lib

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"regexp"
)

type Utils struct {
	Discord *discordgo.Session

	Config *Config

	Log zerolog.Logger

	DB *gorm.DB
}

func (u *Utils) Reply(replyTo *discordgo.Message, message string) {
	_, err := u.Discord.ChannelMessageSendReply(replyTo.ChannelID, message, replyTo.Reference())
	if err != nil {
		u.Log.Error().Msgf("failed to send reply: %s", err)
	}
}

var snowflakeRegexp = regexp.MustCompile(`^\d+$`)

func (u *Utils) UserFromId(userId string) (*discordgo.User, error) {
	if !snowflakeRegexp.MatchString(userId) {
		return nil, errors.New("id isn't a snowflake")
	}

	user, errId := u.Discord.User(userId)
	return user, errId
}

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

	return nil, errors.New("no user with that name is in the server where the message was sent")
}
