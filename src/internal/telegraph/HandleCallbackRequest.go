package telegraph

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	FilesystemWalkRequest        = "FilesystemWalkRequest"
	ListFilesRequest             = "ListFilesRequest"
	FilesystemTextSummaryRequest = "FilesystemTextSummaryRequest"
	DownloadFileRequest          = "DownloadFileRequest"
)

func (t Telegraphist) HandleCallbackRequest(update tgbotapi.Update) {
	user := t.authenticatedUsers[update.CallbackQuery.From.ID]
	if user == nil || !user.authenticated {
		_, _ = t.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You are not authorized, use /authorize <code> to authorize yourself"))
		return
	}
	command, data, err := t.callbackStack.DecodeCallbackRequest(update.CallbackQuery.Data)
	if err != nil {
		log.Println(err)
		t.ReportError(fmt.Sprintf("Error: %v", err), update.CallbackQuery.Message.Chat.ID)
		return
	}
	fmt.Printf("Command: %q | Data: %q\n", command, data)
	switch command {
	case FilesystemWalkRequest:
		err := user.SetPath(data)
		if err != nil {
			t.ReportError("An error occurred while processing filesystemPath request", update.CallbackQuery.Message.Chat.ID)
		}
		dir, err := user.ScanCurrentPath()
		if err != nil {
			t.ReportError(fmt.Sprintf("Error while scanning current path: %v", err), update.CallbackQuery.Message.Chat.ID)
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Directories")
		msg.ReplyMarkup = t.PrepareDirectoriesKeyboard(dir)
		_, err = t.bot.Send(msg)
		if err != nil {
			log.Println(err)
			t.ReportError(fmt.Sprintf("An error occured during sending a message: %v", err), update.CallbackQuery.Message.Chat.ID)
		}
		t.answerCallback(update)
	case FilesystemTextSummaryRequest:
		dir, err := user.ScanPath(data)
		if err != nil {
			log.Println(err)
			t.ReportError(fmt.Sprintf("Error while scanning path: %v", err), update.CallbackQuery.Message.Chat.ID)
			t.answerCallback(update)
			return
		}
		strDir := dir.String()
		if len(strDir) > 4095 {
			doc := tgbotapi.NewDocumentUpload(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileBytes{
				Name:  dir.info.Name(),
				Bytes: []byte(strDir),
			})
			_, err := t.bot.Send(doc)
			if err != nil {
				log.Println(err)
				t.ReportError(fmt.Sprintf("Couldn't send file: %v", err), update.CallbackQuery.Message.Chat.ID)
			}
		} else {
			_, err := t.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, strDir))
			if err != nil {
				log.Println(err)
			}
		}
		t.answerCallback(update)
	case ListFilesRequest:
		dir, err := user.ScanPath(data)
		if err != nil {
			log.Println(err)
			t.ReportError(fmt.Sprintf("Error while scanning path: %v", err), update.CallbackQuery.Message.Chat.ID)
			t.answerCallback(update)
			return
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Files in %v", dir.info.Name()))
		msg.ReplyMarkup = t.PrepareFilesKeyboard(dir)
		_, err = t.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		t.answerCallback(update)
	case DownloadFileRequest:
		fileUpload := tgbotapi.NewDocumentUpload(update.CallbackQuery.Message.Chat.ID, data)
		_, err = t.bot.Send(fileUpload)
		if err != nil {
			log.Println(err)
			t.ReportError(fmt.Sprintf("Couldn't send file: %v", err), update.CallbackQuery.Message.Chat.ID)
		}
		t.answerCallback(update)
	default:
		t.ReportError(fmt.Sprintf("Unknown callback command %v", command), update.CallbackQuery.Message.Chat.ID)
	}
}

func (t Telegraphist) answerCallback(update tgbotapi.Update) {
	_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		log.Println(err)
	}
}
