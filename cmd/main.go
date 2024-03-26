package main

import (
	"log"
	"wordeeBot/internal/bot"
)

func main() {
	bot, err := bot.NewBot("BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	bot.ListenForUpdates()
}
