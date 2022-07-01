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
	// SLHeatMatchupFile is the name of the file where the SL heat matchups are stored
	SLHeatMatchupFile = "/data/images/slMatchupAverages.svg"
	// MatchupSheet is the name of the sheet where the matchups are stored
	MatchupSheet = "Sheet1"
	// ReactionRequest is the reaction emoji choices for availability
	ReactionRequest = "React to this message with a 👍 if you are coming, " +
		"a 👎 if you can't make it, and an ⌛ if you will be late. " +
		"Any reaction of this type in #game-night " +
		"until 7pm will update your tracked availability."
	// DevChannelID is the ID of channel #bot-dev
	DevChannelID = "955291440643207228"
	// GameNightChannelID is the ID of channel #game-night
	GameNightChannelID = "951345352030691381"
	// StrategyChannelID is the ID of channel #strategy
	StrategyChannelID = "951346668912136192"
)

var (
	DivisionTeamNames = [9]string{
		"Jiffyloob",
		"A Selected Few",
		"Wookie Mistakes",
		"In It 2 Win It",
		"The Wright Stuff-8",
		"G Team",
		"The Unusual Suspects",
		"Lil's Bunch",
		"8-Balls of Fire",
	}
	GameDayReactions = []string{"👍", "👎", "⌛", "⏳"}
	SeniorSkillLevel = 6
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
	// MatchupExpectedPoints for the team lineup
	MatchupExpectedPointsFor float64
	// MatchupExpectedPointsAgainst for the team lineup
	MatchupExpectedPointsAgainst float64
	// Matchups for the team lineup
	Matchups []string
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
	// HandleSLTextMatchups for returning chart of the best skill level match-ups
	HandleSLTextMatchups(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleSLHeatMatchups for returning heatmap of the best skill level match-ups
	HandleSLHeatMatchups(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleOptimal for returning max expected points lineup from opponent's lineup
	HandleOptimal(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleBCA for mentions of BCA play
	HandleBCA(s *discordgo.Session, m *discordgo.MessageCreate)
	// Handle9Ball for mentions of 9 ball play
	Handle9Ball(s *discordgo.Session, m *discordgo.MessageCreate)
}
