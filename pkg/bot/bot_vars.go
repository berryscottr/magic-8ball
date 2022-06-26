package bot

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/bwmarrin/discordgo"
)

const (
	// UserID assigns ID to the bot
	UserID = "@me"
	// SLMatchupFile is the name of the file where the SL matchups are stored
	SLMatchupFile = "/data/SLMatchupAverages.xlsx"
	// MatchupSheet is the name of the sheet where the matchups are stored
	MatchupSheet = "Sheet1"
	// ReactionRequest is the reaction emoji choices for availability
	ReactionRequest = "React to this message with a 👍 if you are coming, " +
		"a 👎 if you can't make it, and an ⌛ if you will be late."
)

var (
	DivisionTeamNames = [9]string{
		"Jiffyloob",
		"A Selected Few",
		"Wookie Mistakes",
		"In It 2 Win It",
		"The Wright Stuff",
		"G Team",
		"The Unusual Suspects",
		"Lil's Bunch",
		"8-Balls of Fire",
	}
)

// Data for the bot to track along a request
type Data struct {
	// Err for error tracking
	Err error
	// User for the bot to track
	User *discordgo.User
	// GoBot for the bot to track
	GoBot *discordgo.Session
	// Token for the bot to track
	Token string
	// Excel for the bot to track
	Excel *excelize.File
	// ExcelRows for the bot to track
	ExcelRows [][]string
	// Dir for the bot to track
	Dir string
}

// TeamLineup data
type TeamLineup struct {
	// Lineup for the team lineup
	Lineup []int
	// Sum for the team lineup
	Sum int
}

// Methods for the bot to use
type Methods interface {
	// SetDir for the bot to use
	SetDir() Data
	// Start the Discord bot listener
	Start()
	// MessageHandler for interpreting which function to launch from message contents
	MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate)
	// ReactionHandler for interpreting how to respond to reactions
	ReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd)
	// HandleGameDayReaction for interpreting how to respond to reactions
	HandleGameDayReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd)
	// HandleGameDay for posting game day message
	HandleGameDay(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleLineups for returning eligible lineups from a provided list of players
	HandleLineups(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleSLMatchups for returning eligible lineups from a provided list of players
	HandleSLMatchups(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleBCA for mentions of non-APA play
	HandleBCA(s *discordgo.Session, m *discordgo.MessageCreate)
}
