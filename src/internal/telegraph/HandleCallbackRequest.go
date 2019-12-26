package telegraph

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	FilesystemPathRequest = "FilesystemPathRequest"
)

func (t Telegraphist) HandleCallbackRequest(update tgbotapi.Update) {
	user := t.authenticatedUsers[update.CallbackQuery.From.ID]
	if user == nil || !user.authenticated {
		t.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You are not authorized, use /authorize <code> to authorize yourself"))
		return
	}
	_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		log.Println(err)
	}
	command, data, err := t.callbackStack.DecodeCallbackRequest(update.CallbackQuery.Data)
	if err != nil {
		log.Println(err)
		t.ReportError(fmt.Sprintf("Error: %v", err), update.CallbackQuery.Message.Chat.ID)
		return
	}
	fmt.Printf("Command: %q | Data: %q\n", command, data)
	switch command {
	case FilesystemPathRequest:
		err := user.SetPath(data)
		if err != nil {
			t.ReportError("An error occurred while processing filesystemPath request", update.CallbackQuery.Message.Chat.ID)
		}
		dir, err := user.ScanCurrentPath()
		if err != nil {
			t.ReportError(fmt.Sprintf("Error while scanning current path: %v", err), update.CallbackQuery.Message.Chat.ID)
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Directories")
		msg.ReplyMarkup = t.PrepareFilesystemKeyboard(dir)
		_, err = t.bot.Send(msg)
		if err != nil {
			log.Println(err)
			t.ReportError(fmt.Sprintf("An error occured during sending a message: %v", err), update.CallbackQuery.Message.Chat.ID)
			return
		}

	default:
		t.ReportError("Unknown callback command", update.CallbackQuery.Message.Chat.ID)
	}
}
