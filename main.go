package main

import (
	"github.com/rs/zerolog/log"
	"magic-8ball/pkg/bot"
)

func main() {
	gobot := new(bot.Data)
	gobot.Config = gobot.Config.ReadConfig()
	if gobot.Config.Err != nil {
		gobot.Err = gobot.Config.Err
		log.Err(gobot.Err).Msg("failed to read bot config")
		return
	}
	gobot.Start()
	<-make(chan struct{})
}
