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
	// SLMatchupFileNine is the name of the file where the SL matchups are stored
	SLMatchupFileNine = "/data/SLMatchupAveragesNine.xlsx"
	// EightSLHeatMatchupAveragesUrl is the name of the file where the SL heat matchups are stored
	EightSLHeatMatchupAveragesUrl = "https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupAverages.svg"
	// NineSLHeatMatchupAveragesUrl is the name of the file where the SL heat matchups are stored
	NineSLHeatMatchupAveragesUrl = "https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupAveragesNine.svg"
	// SLMatchupMediansUrl is the name of the file where the SL matchup medians are stored
	SLMatchupMediansUrl = "https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupMedians.png"
	// SLMatchupModesUrl is the name of the file where the SL matchup medians are stored
	SLMatchupModesUrl = "https://raw.githubusercontent.com/berryscottr/magic-8ball/main/data/images/slMatchupModes.png"
	// InningsFile is the name of the file where the SL innings are stored
	InningsFile = "/data/InningCounts.xlsx"
	// ReactionRequest is the reaction emoji choices for availability
	ReactionRequest = "React to this message with a 👍 if you are coming, " +
		"a 👎 if you can't make it, and an ⌛ if you will be late. " +
		"Any reaction of this type in a #game-night channel " +
		"until 7pm will send an update message to the #game-night channels."
	// DevChannelID is the ID of channel #bot-dev
	DevChannelID = "955291440643207228"
	// GameNight8ChannelID is the ID of channel #game-night
	GameNight8ChannelID = "951345352030691381"
	// GameNight9ChannelID is the ID of channel #game-night
	GameNight9ChannelID = "1013889839101399111"
	// EightBallRoleID is the ID of role 8-Ball
	EightBallRoleID = "1013886913880522872"
	// NineBallRoleID is the ID of role 9-Ball
	NineBallRoleID = "1013887160480436317"
	// StrategyChannelID is the ID of channel #strategy
	StrategyChannelID = "951346668912136192"
	// APACalendarUrl is the URL of the APA calendar
	APACalendarUrl = "https://atlanta.apaleagues.com/Uploads/atlanta/APA%20Atlanta%202023.pdf"
	// TeamCalendarUrl is the URL of the team calendar
	TeamCalendarUrl = "https://github.com/berryscottr/magic-8ball/blob/main/data/schedules/Spring2023Schedule.csv"
	// SeniorSkillLevel is the skill level of seniors
	SeniorSkillLevel = 6
)

type WookieMistakes struct {
	// TeamName for the team
	TeamName string
	// DivisionTeamNames for the team
	DivisionTeamNames [8]string
}

type SafetyDance struct {
	// TeamName for the team
	TeamName string
	// DivisionTeamNames for the team
	DivisionTeamNames [7]string
}


var (
	WookieMistakesTeam = WookieMistakes{
		TeamName: "Wookie Mistakes",
		DivisionTeamNames: [8]string{
			"Shark Shooters - 8",
			"Only A Few Selected",
			"In It 2 Win It",
			"Jiffyloob",
			"G Team",
			"The Unusual Suspects",
			"School of Pool",
			"8-Balls of Fire",
		},
	}
	SafetyDanceTeam = SafetyDance{
		TeamName: "Safety Dance",
		DivisionTeamNames: [7]string{
			"Shark Shooters - 9",
			"Sticks and Stones",
			"9 Rocks Away",
			"9 on the Vine",
			"Rackateers",
			"Believe It or Not 2",
			"Fields of Gold 9",
		},
	}
	// GameDayReactions for the bot to track
	GameDayReactions = []string{"👍", "👎", "⌛", "⏳"}
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
	// OpponentLineup for the team lineup
	OpponentLineup []int
	// Sum for the team lineup
	Sum int
	// MatchupExpectedPointsFor the team lineup
	MatchupExpectedPointsFor float64
	// MatchupExpectedPointsAgainst for the team lineup
	MatchupExpectedPointsAgainst float64
	// MatchupExpectedPointsDifference ExpectedPointsFor - ExpectedPointsAgainst
	MatchupExpectedPointsDifference float64
	// Matchups for the team lineup
	Matchups []Matchup
}

type Matchup struct {
	// SkillLevels for the match-up
	SkillLevels [2]int
	// ExpectedPointsFor for the match-up
	ExpectedPointsFor float64
	// ExpectedPointsAgainst for the match-up
	ExpectedPointsAgainst float64
	// ExpectedPointsDifference ExpectedPointsFor - ExpectedPointsAgainst
	ExpectedPointsDifference float64
}

// Methods for the bot to use
type Methods interface {
	// SetDir for the bot to use
	SetDir()
	// Start the Discord bot listener
	Start()
	// MessageHandler for interpreting which function to launch from message contents
	MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate)
	// ReactionHandler for interpreting how to respond to reactions
	ReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd)
	// HandleGameDayReaction8 for interpreting how to respond to reactions
	HandleGameDayReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd)
	// HandleGameDay8 for posting game day message
	HandleGameDay8(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleGameDay9 for posting game day message
	HandleGameDay9(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleLineups for returning eligible lineups from a provided list of players
	HandleLineups(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleSLMatchups for returning chart of the best skill level match-ups
	HandleSLMatchups(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleHandicapAvg for returning your effective innings per game
	HandleHandicapAvg(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleOptimal8 for returning max expected points lineup from opponent's lineup for eight-ball
	HandleOptimal8(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleOptimal9 for returning max expected points lineup from opponent's lineup for nine-ball
	HandleOptimal9(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandlePlayoff for returning max differential expected points lineup from opponent's lineup
	HandlePlayoff(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleCalendar for returning the current calendar
	HandleCalendar(s *discordgo.Session, m *discordgo.MessageCreate)
	// HandleBCA for mentions of BCA play
	HandleBCA(s *discordgo.Session, m *discordgo.MessageCreate)
	// Handle9Ball for mentions of 9 ball play
	Handle9Ball(s *discordgo.Session, m *discordgo.MessageCreate)
}
