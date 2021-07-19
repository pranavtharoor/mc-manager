package main

import (
	"log"

	"github.com/pranavtharoor/mc-manager/bot"
	"github.com/pranavtharoor/mc-manager/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if err := bot.Start(conf.Bot); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	log.Println("Bot running")

	<-make(chan struct{})
}
