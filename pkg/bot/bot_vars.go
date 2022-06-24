package bot

import (
	"github.com/bwmarrin/discordgo"
)

// UserID assigns ID to the bot
const UserID = "@me"

// Data for the bot to track along a request
type Data struct {
	Err   error
	User  *discordgo.User
	GoBot *discordgo.Session
	Token string
}

// TeamLineup data
type TeamLineup struct {
	Lineup []int
	Sum    int
}

// Methods for the bot to use
type Methods interface {
	Start()
	MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate)
	HandleLineups(s *discordgo.Session, m *discordgo.MessageCreate)
}
