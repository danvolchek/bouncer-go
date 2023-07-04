package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/danvolchek/bouncer-go/database"
	"github.com/danvolchek/bouncer-go/lib"
	"github.com/danvolchek/bouncer-go/lib/components"
	"github.com/rs/zerolog/log"
)

type help struct{}

func (h help) Name() string {
	return "help"
}

func (h help) Handle(command *components.CommandDetails, message *discordgo.Message, utils *lib.Utils) bool {
	utils.Log.Info().Msg("this is a random info message")

	//utils.Reply(message, "hi!")

	var thread database.UserReplyThread

	if err := utils.DB.First(&thread).Error; err != nil {
		utils.Log.Error().Err(err).Msg("no user waaah")
		return false
	}

	log.Info().Msg(fmt.Sprintf("%+v", thread))

	return false
}
