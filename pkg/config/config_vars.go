package config

var FilePath = "config.json"

// Conf contains data for reading bot config
type Conf struct {
	Err      error
	FilePath string
	File     []byte
	Data     File
}

// File is generated from the config json
type File struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

// Methods of the config package
type Methods interface {
	ReadConfig() Conf
}
