package main

import (
	"magic-8ball/pkg/bot"
	"magic-8ball/pkg/healthz"
	"os"
	"os/signal"
	"syscall"
	"github.com/rs/zerolog/log"
	"github.com/bwmarrin/discordgo"
)

func main() {
	gobot := bot.Data{Token: bot.Token{Discord: os.Getenv("BOT_TOKEN")}}
	if gobot.Token.Discord == "" {
		panic("BOT_TOKEN not set")
	}
	gobot.SetDir()
	gobot.Start()
	health := healthz.HealthCheckServer{}
	health.Start()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	session, err := discordgo.New("Bot " + gobot.Token.Discord)
	if err != nil {
			log.Fatal().Err(err).Msg("failed to create Discord session")
	}
	go gobot.ScheduleGameDay(session, nil, "Wookie Mistakes")
	go gobot.ScheduleGameDay(session, nil, "Safety Dance")

	<-sigs

	log.Info().Msg("shutting down magic-8ball bot")
	gobot.Stop()
	health.Close()
}
