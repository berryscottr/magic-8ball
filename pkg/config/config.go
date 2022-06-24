package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

func (config Conf) ReadConfig() Conf {
	config.File, config.Err = ioutil.ReadFile("./config.json")
	if config.Err != nil {
		log.Err(config.Err).Msg("failed to read config file")
		return config
	}
	config.Err = json.Unmarshal(config.File, &config.Data)
	if config.Err != nil {
		log.Err(config.Err).Msg("failed to unmarshal config bytes")
		return config
	}
	return config
}
