package main

import (
	"magic-8ball/pkg/bot"
)

func main() {
	gobot := new(bot.Data)
	gobot.Start()
	<-make(chan struct{})
}
