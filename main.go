package main

import (
	"magic-8ball/pkg/bot"
	"os"
)

func main() {
	gobot := bot.Data{Token: os.Getenv("BOT_TOKEN")}
	gobot.Start()
	<-make(chan struct{})
}
