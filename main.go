package main

import (
	"magic-8ball/pkg/bot"
	"magic-8ball/pkg/healthz"
	"os"
	"os/signal"
	"syscall"
	"github.com/rs/zerolog/log"
)

func main() {
	gobot := bot.Data{Token: os.Getenv("BOT_TOKEN")}
	if gobot.Token == "" {
		panic("BOT_TOKEN not set")
	}
	gobot.SetDir()
	gobot.Start()
	health := healthz.HealthCheckServer{}
	health.Start()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info().Msg("shutting down magic-8ball bot")
	gobot.Stop()
	health.Close()
}
