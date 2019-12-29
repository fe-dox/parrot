package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(settings.BotToken)
	if err != nil {
		log.Panic(err)
	}
	telegraphist := NewTelegraphist(bot)
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			go telegraphist.HandleCallbackRequest(update)
		}
		if update.Message == nil {
			continue
		}
		go telegraphist.HandleCommand(update)

	}
}
