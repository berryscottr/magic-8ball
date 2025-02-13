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
	ReactionRequest = "React to this message with the following choices:\n👍 Available \t👎 Unavailable\n⌛ Late      \t\t❓ Unsure"
	// DevChannelID is the ID of channel #bot-dev
	DevChannelID = "955291440643207228"
	// TestChannelID is the ID of channel #bot-testing
	TestChannelID = "1201198112564318288"
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

type Team struct {
	// Format for the team
	Format string
	// Name for the team
	Name string
	// DivisionTeamNames for the team
	DivisionTeamNames []string
	// GameDay for the team
	GameDay string
	// GameNightChannelID for the team
	GameNightChannelID string
}

// Teammate Name and ID
type Teammate struct {
	// LastName
	LastName string
	// FirstName
	FirstName string
	// UserID
	UserID string
	// Name
	Name string
	// Teams
	Teams []Team
	// Skill Level
	SkillLevel Skill
	// Roster Number
	RosterNum Roster
}

// Skill Level
type Skill struct {
	// 8-Ball
	Eight int
	// 9-Ball
	Nine int
}

// Roster number
type Roster struct {
	// 8-Ball
	Eight int
	// 9-Ball
	Nine int
}

// TeammateSkill for a teammate and their skill level
type TeammateSkill struct {
	// Teammate
	Lastname string
	// Skill
	SkillLevel int
}

var (
	// Teammates on teams
	Teammates = []Teammate{
		{
			LastName: "Berry",
			FirstName: "Scott",
			UserID:   "341590471317127178",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 6,
				Nine:  6,
			},
			RosterNum: Roster{
				Eight: 1,
				Nine:  1,
			},
		},
		{
			LastName: "Liess",
			FirstName: "Jason",
			UserID:   "529730365854580765",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 5,
				Nine:  4,
			},
			RosterNum: Roster{
				Eight: 2,
				Nine:  2,
			},
		},
		{
			LastName: "Bohrer",
			UserID:   "186997536844808193",
			Teams:    []Team{},
			SkillLevel: Skill{
				Eight: 3,
				Nine:  3,
			},
		},
		{
			LastName: "Burcham",
			FirstName: "Daniel",
			UserID:   "1014488206567288894",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 5,
				Nine:  4,
			},
			RosterNum: Roster{
				Eight: 3,
				Nine:  3,
			},
		},
		{
			LastName: "Thompson",
			FirstName: "Laura",
			UserID:   "969682397920653342",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 3,
				Nine:  3,
			},
			RosterNum: Roster{
				Eight: 4,
				Nine:  4,
			},
		},
		{
			LastName: "Quan",
			FirstName: "Kim",
			UserID:   "795533691828305922",
			Teams: []Team{
				WookieMistakes,
			},
			SkillLevel: Skill{
				Eight: 5,
			},
			RosterNum: Roster{
				Eight: 5,
			},
		},
		{
			LastName: "Hayward",
			FirstName: "Leo",
			UserID:   "971791774697783326",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 7,
				Nine:  8,
			},
			RosterNum: Roster{
				Eight: 6,
				Nine:  5,
			},
		},
		{
			LastName: "Davalos",
			UserID:   "1108221800581709917",
			Teams:    []Team{},
			SkillLevel: Skill{
				Nine: 1,
			},
		},
		{
			LastName: "Warden",
			UserID:   "1014520790873546852",
			Teams:    []Team{},
			SkillLevel: Skill{
				Eight: 7,
				Nine:  9,
			},
		},
		{
			LastName: "Gibson",
			FirstName: "Alex",
			UserID:   "696037354892296294",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 4,
				Nine:  3,
			},
			RosterNum: Roster{
				Eight: 7,
				Nine:  6,
			},
		},
		{
			LastName: "Dodge",
			FirstName: "Dylan",
			UserID:   "253692229535793154",
			Teams: []Team{
				WookieMistakes,
				SafetyDance,
			},
			SkillLevel: Skill{
				Eight: 4,
				Nine:  4,
			},
			RosterNum: Roster{
				Eight: 8,
				Nine:  7,
			},
		},
		{
			LastName: "Ackler",
			UserID:   "969934577768484944",
			Teams: []Team{},
			SkillLevel: Skill{
				Eight: 6,
				Nine:  6,
			},
		},
		{
			LastName: "Johnson",
			UserID:   "1236395641215651850",
			Teams: []Team{},
			SkillLevel: Skill{
				Eight: 3,
				Nine:  4,
			},
			RosterNum: Roster{},
		},
		{
			LastName: "Gunning",
			FirstName: "Matt",
			UserID:   "360240801353302016",
			Teams: []Team{
				SafetyDance,
			},
			SkillLevel: Skill{
				Nine: 3,
			},
			RosterNum: Roster{
				Nine: 8,
			},
		},
	}

	// WookieMistakes Eight-Ball team
	WookieMistakes = Team{
		Format: "8-Ball",
		Name:   "Wookie Mistakes",
		DivisionTeamNames: []string{
			"Shark Shooters - 8",
			"Jiffyloob",
			"The Unusual Suspects",
			"School of Pool",
			"Ex-Cues US",
		},
		GameDay:            "Tuesday",
		GameNightChannelID: GameNight8ChannelID,
	}
	// SafetyDance Nine-Ball team
	SafetyDance = Team{
		Format: "9-Ball",
		Name:   "Safety Dance",
		DivisionTeamNames: []string{
			"Shark Shooters - 9",
			"9 Rocks Away",
			"9 on the Vine",
			"In The Pocket-9",
			"Believe It or Not 2",
			"Salted Coffee",
			"Smooth Strokes-9",
		},
		GameDay:            "Tuesday",
		GameNightChannelID: GameNight9ChannelID,
	}
	// GameDayReactions for the bot to track
	GameDayReactions = []string{"👍", "👎", "⌛", "⏳", "❓", "❔", NumToEmojiMap[1], NumToEmojiMap[2], NumToEmojiMap[3], NumToEmojiMap[4], NumToEmojiMap[5], NumToEmojiMap[6], NumToEmojiMap[7], NumToEmojiMap[8]}
	// numToEmojiMap is a map for converting numbers to emojis
	NumToEmojiMap = map[int]string{
		1: "1️⃣",
		2: "2️⃣",
		3: "3️⃣",
		4: "4️⃣",
		5: "5️⃣",
		6: "6️⃣",
		7: "7️⃣",
		8: "8️⃣",
		9: "9️⃣",
	}
	EightBallEmoji = "🎱"
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
	// Stop the Discord bot listener
	Stop()
	// MessageHandler for interpreting which function to launch from message contents
	MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate)
	// ReactionHandler for interpreting how to respond to reactions
	ReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd)
	// HandleGameDayReaction8 for interpreting how to respond to reactions
	HandleGameDayReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd)
	// HandleGameDay for posting game day message
	HandleGameDay(s *discordgo.Session, m *discordgo.MessageCreate, teamName string)
	// ScheduleGameDay for posting game day message
	ScheduleGameDay(s *discordgo.Session, m *discordgo.MessageCreate, teamName string)
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
}

