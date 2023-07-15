package components

import (
	"github.com/bwmarrin/discordgo"
	"github.com/danvolchek/bouncer-go/lib"
)

type Ready struct {
	*lib.Utils
}

func NewReady() *Ready {
	return &Ready{}
}

func (r *Ready) Setup(utils *lib.Utils) {
	r.Utils = utils

	r.Discord.AddHandler(r.ready)
}

func (r *Ready) ready(_ *discordgo.Session, _ *discordgo.Ready) {
	r.Log.Info().Msg("Connected to discord!")
}
