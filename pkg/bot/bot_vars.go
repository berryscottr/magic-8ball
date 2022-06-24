package bot

import (
	"github.com/bwmarrin/discordgo"
	"magic-8ball/pkg/config"
)

const UserID = "@me"

type Data struct {
	Err    error
	User   *discordgo.User
	GoBot  *discordgo.Session
	Config config.Conf
}

type TeamLineup struct {
	Lineup []int
	Sum    int
}

type Methods interface {
	Start()
	MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate)
	HandleLineups(s *discordgo.Session, m *discordgo.MessageCreate)
}
