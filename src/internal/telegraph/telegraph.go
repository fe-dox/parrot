package telegraph

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Telegraph interface {
	HandleCommand(update tgbotapi.Update)
}

type Telegraphist struct {
	bot                *tgbotapi.BotAPI
	authenticatedUsers map[int]User
}

func NewTelegraphist(bot *tgbotapi.BotAPI) *Telegraphist {
	return &Telegraphist{bot: bot, authenticatedUsers: make(map[int]User)}
}
