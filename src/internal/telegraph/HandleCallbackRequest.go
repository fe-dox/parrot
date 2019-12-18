package telegraph

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

const (
	FilesystemPathRequest = "FilesystemPathRequest"
)

func (t Telegraphist) HandleCallbackRequest(update tgbotapi.Update) {
	if !t.authenticatedUsers[update.Message.From.ID].authenticated {
		t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized, use /authorize <code> to authorize yourself"))
		return
	}
	_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		log.Println(err)
	}
	innerUser := t.authenticatedUsers[update.CallbackQuery.From.ID]
	command, data := splitCommandAndData(update.CallbackQuery.Data)

	switch command {
	case FilesystemPathRequest:
		err := innerUser.SetPath(data)
		if err != nil {
			t.ReportError("An error occurred while processing filesystemPath request", update.Message.Chat.ID)
		}

	}
}

func splitCommandAndData(s string) (string, string) {
	i := strings.Index(s, "-")
	return s[:i], s[i+1:]
}
