package main

import (
	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vsergeev/btckeygenie/server"
)

func main() {
	c := &server.Config{}
	err := env.Parse(c)
	if err != nil {
		log.Fatal().Err(err).Msg("env.Parse failed")
	}
	setLogLevel(c)
	log.Info().Str("config", c.String()).Send()

	s, err := c.InitServer()
	if err != nil {
		log.Fatal().Err(err).Msg("InitServer failed")
	}
	defer s.DB.Close()

	err = s.Serve()
	if err != nil {
		log.Fatal().Err(err).Msg("Run failed")
	}
}

func setLogLevel(cfg *server.Config) {
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug is enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
