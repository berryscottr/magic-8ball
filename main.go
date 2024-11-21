package main

import (
	// "magic-8ball/pkg/bot"
	"fmt"
	"magic-8ball/pkg/league"
	// "os"
)

func main() {
	league := league.Data{}
	err := league.Login()
	if err != nil {
		fmt.Print(err)
	}
// 	gobot := bot.Data{Token: os.Getenv("BOT_TOKEN")}
// 	if gobot.Token == "" {
// 		panic("BOT_TOKEN not set")
// 	}
// 	gobot.SetDir()
// 	gobot.Start()
// 	<-make(chan struct{})
}
