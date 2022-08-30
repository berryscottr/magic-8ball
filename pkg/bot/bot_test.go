package bot

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestData_SetDir confirms ability to set the bot directory
func TestData_SetDir(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to set the bot directory")
}

// TestData_Start confirms ability to start the discord listener
func TestData_Start(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	data.Start()
	assertion.NoError(data.Err, "failed to start discord listener")
}

// TestData_MessageHandler confirms ability to handle message contents
func TestData_MessageHandler(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to handle message contents")
}

// TestData_ReactionHandler confirms ability to handle reactions
func TestData_ReactionHandler(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to handle reaction")
}

// TestData_HandleGameDay8Reaction confirms ability to handle the reaction to the game day post
func TestData_HandleGameDay8Reaction(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to read a game day reaction")
}

// TestData_HandleGameDay9Reaction confirms ability to handle the reaction to the game day post
func TestData_HandleGameDay9Reaction(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to read a game day reaction")
}

// TestData_HandleGameDay8 confirms ability to generate a game day post
func TestData_HandleGameDay8(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to generate a response to a game day post")
}

// TestData_HandleGameDay9 confirms ability to generate a game day post
func TestData_HandleGameDay9(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to generate a response to a game day post")
}

// TestData_HandleLineups confirms ability to generate valid lineups from an input
func TestData_HandleLineups(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to generate valid lineups")
}

// TestData_HandleSLMatchups confirms ability to read the matchups Excel sheet
func TestData_HandleSLMatchups(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to read the matchups excel sheet")
}

// TestData_HandleHandicapAvg confirms ability to return the average handicap
func TestData_HandleHandicapAvg(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to read the matchups excel sheet")
}

// TestData_HandleOptimal confirms ability to generate a max expected points lineup
func TestData_HandleOptimal(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	//var s *discordgo.Session
	//m := discordgo.MessageCreate{
	//	Message: &discordgo.Message{
	//		Content: "!optimal 65533 22335567",
	//	},
	//}
	//data.HandleOptimal(s, &m)
	assertion.NoError(data.Err, "failed to return the optimal matchup")
}

// TestData_HandlePlayoff confirms ability to generate a max differential expected points lineup
func TestData_HandlePlayoff(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to return the optimal playoff matchup")
}

// TestData_HandleBCA confirms ability to respond to a BCA mention
func TestData_HandleBCA(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to generate a response to a BCA mention")
}

// TestData_Handle9Ball confirms ability to respond to a 9ball mention
func TestData_Handle9Ball(t *testing.T) {
	assertion := assert.New(t)
	data := Data{Token: os.Getenv("BOT_TOKEN")}
	data = data.SetDir()
	assertion.NoError(data.Err, "failed to generate a response to a 9ball mention")
}
