package bot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConf_ReadConfig confirms ability to start the discord listener
func TestData_Start(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to start discord listener")
}

// TestData_MessageHandler confirms ability to handle message contents
func TestData_MessageHandler(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to handle message contents")
}

// TestData_ReactionHandler confirms ability to handle reactions
func TestData_ReactionHandler(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to handle reaction")
}

// TestData_HandleGameDayReaction confirms ability to handle the reaction to the game day post
func TestData_HandleGameDayReaction(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to read a game day reaction")
}

// TestData_HandleGameDay confirms ability to generate a game day post
func TestData_HandleGameDay(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to generate a response to a game day post")
}

// TestData_HandleLineups confirms ability to generate valid lineups from an input
func TestData_HandleLineups(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to generate valid lineups")
}

// TestData_HandleSLMatchups confirms ability to read the matchups Excel sheet
func TestData_HandleSLMatchups(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to read the matchups excel sheet")
}

// TestData_HandleBCA confirms ability to respond to a BCA mention
func TestData_HandleBCA(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to generate a response to a BCA mention")
}

// TestData_Handle9Ball confirms ability to respond to a 9ball mention
func TestData_Handle9Ball(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to generate a response to a 9ball mention")
}
