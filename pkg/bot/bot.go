package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

// SetDir for setting the directory for the bot
func (bot Data) SetDir() Data {
	bot.Dir, bot.Err = filepath.Abs(".")
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to set magic-8ball directory")
	}
	return bot
}

// Start the Discord bot listener
func (bot Data) Start() {
	bot.GoBot, bot.Err = discordgo.New("Bot " + bot.Token)
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to instantiate magic-8ball bot")
		return
	}
	bot.User, bot.Err = bot.GoBot.User(UserID)
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to set magic-8ball user id")
		return
	}
	bot.GoBot.AddHandler(bot.MessageHandler)
	bot.GoBot.AddHandler(bot.ReactionHandler)
	bot.Err = bot.GoBot.Open()
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to start magic-8ball listener")
		return
	}
	log.Info().Msg("magic-8ball listening")
}
