package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

// SetDir for setting the directory for the bot
func (bot *Data) SetDir() {
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
