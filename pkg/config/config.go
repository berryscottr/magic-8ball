package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path/filepath"
)

// ReadConfig from the config file
func (config Conf) ReadConfig() Conf {
	config.FilePath, config.Err = filepath.Abs(FilePath)
	if config.Err != nil {
		log.Err(config.Err).Msg("failed to find config file")
		return config
	}
	config.File, config.Err = ioutil.ReadFile(config.FilePath)
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
