package main

import (
	"encoding/json"
	"flag"
	"github.com/danvolchek/bouncer-go/commands"
	"github.com/danvolchek/bouncer-go/database"
	"github.com/danvolchek/bouncer-go/lib"
	"github.com/danvolchek/bouncer-go/lib/components"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	configPath := flag.String("config", "", "path to config file (required)")
	debug := flag.Bool("debug", false, "sets log level to debug")
	trace := flag.Bool("trace", false, "sets log level to trace")

	// parse args
	{
		flag.Parse()

		if *configPath == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	// set up logging
	{
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateOnly + " " + time.Kitchen + " MST"}).
			With().Caller().
			Logger()

		logLevel := zerolog.InfoLevel
		if *trace {
			logLevel = zerolog.TraceLevel
		} else if *debug {
			logLevel = zerolog.DebugLevel
		}

		zerolog.SetGlobalLevel(logLevel)
		log.Info().Msgf("Log level is %s", logLevel)
	}

	// Load config
	var config *lib.Config
	{
		data, err := os.ReadFile(*configPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read config file")
		}

		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse config file")
		}
	}

	// Initialize database
	db, err := database.New(config.DatabasePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database")
	}

	// create components
	var comps []lib.Component
	{
		commandHandler, err := components.NewCommandHandler(commands.All)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create command handler")
		}

		readyLogger := components.NewReadyLogger()

		comps = []lib.Component{commandHandler, readyLogger}
	}

	// create and run bot
	{
		bot := lib.NewBot(comps, config, db)

		err = bot.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to run bot")
		}
	}
}
