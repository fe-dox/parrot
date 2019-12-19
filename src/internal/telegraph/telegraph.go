package telegraph

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Telegraph interface {
	HandleCommand(update tgbotapi.Update)
	HandleCallbackRequest(update tgbotapi.Update)
	PrepareFilesystemKeyboard(d Directory) tgbotapi.InlineKeyboardMarkup
}

type Telegraphist struct {
	bot                *tgbotapi.BotAPI
	authenticatedUsers map[int]User
	callbackStack      CallbackStack
}

type (
	CallbackStack       map[CallbackID]map[CallbackStackItemID]CallbackStackItem
	CallbackID          int
	CallbackStackItemID int
	CallbackStackItem   struct {
		Command string
		Data    string
	}
)

type StackResolver interface {
	AddCallback(seed int64) CallbackID
	AddCallbackItem(id CallbackID) CallbackStackItemID
}

func NewTelegraphist(bot *tgbotapi.BotAPI) *Telegraphist {
	return &Telegraphist{bot: bot, authenticatedUsers: make(map[int]User), callbackStack: make(CallbackStack)}
}
