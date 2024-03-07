package bot

import (
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// SetDir for setting the directory for the bot
func (bot *Data) SetDir() {
	if _, bot.Err = os.Stat("/app"); bot.Err != nil && os.IsNotExist(bot.Err) {
			bot.Dir, bot.Err = os.Getwd()
			if bot.Err != nil {
					log.Err(bot.Err).Msg("failed to get current directory")
					return
			}
			return
	}
	if bot.Err = os.Chdir("/app"); bot.Err != nil {
			log.Err(bot.Err).Msg("failed to change directory to /app")
			bot.Dir = "/"
			return
	}
	if bot.Dir, bot.Err = filepath.Abs("."); bot.Err != nil {
			log.Err(bot.Err).Msg("failed to set magic-8ball directory")
	}
}

// Start the Discord bot listener
func (bot *Data) Start() {
	if bot.GoBot, bot.Err = discordgo.New("Bot " + bot.Token); bot.Err != nil {
		log.Err(bot.Err).Msg("failed to instantiate magic-8ball bot")
		return
	}
	if bot.User, bot.Err = bot.GoBot.User(UserID); bot.Err != nil {
		log.Err(bot.Err).Msg("failed to set magic-8ball user id")
		return
	}
	bot.GoBot.AddHandler(bot.MessageHandler)
	bot.GoBot.AddHandler(bot.ReactionHandler)
	if bot.Err = bot.GoBot.Open(); bot.Err != nil {
		log.Err(bot.Err).Msg("failed to start magic-8ball listener")
		return
	}
	log.Info().Msg("magic-8ball listening")
}
