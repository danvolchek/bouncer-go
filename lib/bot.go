package lib

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

// Bot is a discord bot. It's composed of a set of components - each which performs a specific bot functionality.
// This struct is the glue that loads the config file, gets all the components ready and connects to discord.
type Bot struct {
	// the components to run
	components []Component

	// general utilities
	*Utils
}

// Component is the interface that bot components should implement.
type Component interface {
	// Setup is called before the bot is started. Register any hooks/perform any initial setup here.
	// The config/db are ready here, but while the discord client has been created it hasn't connected to discord yet,
	// so do not make discord API calls.
	Setup(utils *Utils)
}

// NewBot creates a new discord bot.
func NewBot(components []Component, config *Config, db *gorm.DB, log zerolog.Logger) *Bot {
	return &Bot{
		components: components,

		Utils: &Utils{
			Config: config,
			Log:    log,
			DB:     db,
		},
	}
}

// Run runs the bot until ctrl+c is entered.
func (b *Bot) Run() error {
	discord, err := discordgo.New("Bot " + b.Config.Token)
	if err != nil {
		return err
	}

	b.Utils.Discord = discord

	for _, component := range b.components {
		component.Setup(b.Utils)
	}

	discord.Identify.Intents = discordgo.IntentsAll

	err = discord.Open()
	if err != nil {
		return err
	}

	b.Log.Info().Msg("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = discord.Close()
	if err != nil {
		b.Log.Warn().Msgf("failed to close discord connection (this may impact future runs): %s", err)
	}

	return nil
}
