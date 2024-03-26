package main

import (
	"log"
	"os"
	"wordeeBot/internal/bot"
)

func main() {
	bot, err := bot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.ListenForUpdates()
}
