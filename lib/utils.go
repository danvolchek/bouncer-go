package lib

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
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
