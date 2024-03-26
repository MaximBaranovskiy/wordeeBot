package main

import (
	"log"
	"wordeeBot/internal/bot"
)

func main() {
	bot, err := bot.NewClient("BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	bot.ListenForUpdates()
}
