package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./fixtures/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Err(err)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if s, err := NewServer(); err != nil {
		log.Err(err)
	} else {
		s.Run()
	}
}
