package main

import (
	"fmt"
	"github.com/berryscottr/pkg/bot"
	"github.com/berryscottr/pkg/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
