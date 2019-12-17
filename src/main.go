package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"internal/settings"
	"internal/telegraph"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(settings.BotToken)
	if err != nil {
		log.Panic(err)
	}
	telegraphist := telegraph.NewTelegraphist(bot)
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		go telegraphist.HandleCommand(update)
	}
}
