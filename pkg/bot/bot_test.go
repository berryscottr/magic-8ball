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

// TestData_HandleLineups confirms ability to generate valid lineups from an input
func TestData_HandleLineups(t *testing.T) {
	assertion := assert.New(t)
	data := new(Data)
	assertion.NoError(data.Err, "failed to generate valid lineups")
}
