package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConf_ReadConfig confirms ability to read config file
func TestConf_ReadConfig(t *testing.T) {
	assertion := assert.New(t)
	conf := new(Conf)
	FilePath = "../../" + FilePath
	*conf = conf.ReadConfig()
	assertion.NoError(conf.Err, "failed to read config")
}
