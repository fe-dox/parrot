package telegraph

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
)

type Telegraph interface {
	HandleCommand(update tgbotapi.Update)
	HandleCallbackRequest(update tgbotapi.Update)
	PrepareFilesystemKeyboard(d Directory) tgbotapi.InlineKeyboardMarkup
}

type Telegraphist struct {
	bot                *tgbotapi.BotAPI
	authenticatedUsers map[int]User
}

func NewTelegraphist(bot *tgbotapi.BotAPI) *Telegraphist {
	return &Telegraphist{bot: bot, authenticatedUsers: make(map[int]User)}
}

func (t Telegraphist) PrepareFilesystemKeyboard(d Directory) tgbotapi.InlineKeyboardMarkup {
	keyboardRow := make([]tgbotapi.InlineKeyboardButton, len(d.innerDirs))
	for i, v := range d.innerDirs {
		keyboardRow[i] = tgbotapi.NewInlineKeyboardButtonData(v.name, fmt.Sprintf("%v-%v", FilesystemPathRequest, d.path))
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboardRow)
}
