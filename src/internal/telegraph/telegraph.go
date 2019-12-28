package telegraph

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Telegraph interface {
	HandleCommand(update tgbotapi.Update)
	HandleCallbackRequest(update tgbotapi.Update)
	PrepareDirectoriesKeyboard(d Directory) tgbotapi.InlineKeyboardMarkup
	PrepareFilesKeyboard(d Directory) tgbotapi.InlineKeyboardMarkup
}

type Telegraphist struct {
	bot                *tgbotapi.BotAPI
	authenticatedUsers map[int]*User
	callbackStack      CallbackStack
}

func NewTelegraphist(bot *tgbotapi.BotAPI) *Telegraphist {
	return &Telegraphist{bot: bot, authenticatedUsers: make(map[int]*User), callbackStack: make(CallbackStack)}
}
