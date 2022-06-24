package config

type Conf struct {
	Err  error
	File []byte
	Data File
}

type File struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

type Methods interface {
	ReadConfig() Conf
}
