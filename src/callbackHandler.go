package main

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
	ChangePathRequest            = "ChangePathRequest"
	SendMessageRequest           = "SendMessageRequest"
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
		t.QuickSend(fmt.Sprintf("Error: %v", err), update.CallbackQuery.Message.Chat.ID)
		return
	}
	fmt.Printf("Command: %q | Data: %q\n", command, data)
	switch command {
	case FilesystemWalkRequest:
		err := user.SetPath(data)
		if err != nil {
			t.QuickSend("An error occurred while processing filesystemPath request", update.CallbackQuery.Message.Chat.ID)
			return
		}
		dir, err := user.ScanCurrentPath()
		if err != nil {
			t.QuickSend(fmt.Sprintf("Error while scanning current path: %v", err), update.CallbackQuery.Message.Chat.ID)
			return
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Directories inside %q\n\nNumber of directories inside: %v\nNumber of files inside: %v", dir.info.Name(), len(dir.innerDirs), len(dir.innerFiles)))
		msg.ReplyMarkup = t.PrepareDirectoriesKeyboard(dir)
		_, err = t.bot.Send(msg)
		if err != nil {
			log.Println(err)
			t.QuickSend(fmt.Sprintf("An error occured during sending a message: %v", err), update.CallbackQuery.Message.Chat.ID)
		}
		t.answerCallback(update)
	case FilesystemTextSummaryRequest:
		dir, err := user.ScanPath(data)
		if err != nil {
			log.Println(err)
			t.QuickSend(fmt.Sprintf("Error while scanning path: %v", err), update.CallbackQuery.Message.Chat.ID)
			t.answerCallback(update)
			return
		}
		strDir := dir.String()
		if len(strDir) > 4095 {
			doc := tgbotapi.NewDocumentUpload(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileBytes{
				Name:  dir.info.Name() + ".txt",
				Bytes: []byte(strDir),
			})
			_, err := t.bot.Send(doc)
			if err != nil {
				log.Println(err)
				t.QuickSend(fmt.Sprintf("Couldn't send file: %v", err), update.CallbackQuery.Message.Chat.ID)
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
			t.QuickSend(fmt.Sprintf("Error while scanning path: %v", err), update.CallbackQuery.Message.Chat.ID)
			t.answerCallback(update)
			return
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Files in %q\n\nNumber of directories: %v\nNumber of files: %v", dir.info.Name(), len(dir.innerDirs), len(dir.innerFiles)))
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
			t.QuickSend(fmt.Sprintf("Couldn't send file: %v", err), update.CallbackQuery.Message.Chat.ID)
		}
		t.answerCallback(update)
	case SendMessageRequest:
		t.QuickSend(data, update.CallbackQuery.Message.Chat.ID)
		t.answerCallback(update)
	case ChangePathRequest:
		csID, _, _ := DecodeString(update.CallbackQuery.Data)
		if csID != 0 {
			t.callbackStack.ClearCallback(csID)
		}
		err := user.SetPath(data)
		if err != nil {
			t.QuickSend(fmt.Sprintf("Couldn't change path to %q", data), update.CallbackQuery.Message.Chat.ID)
			t.answerCallback(update)
			return
		}
		t.QuickSend(fmt.Sprintf("Path changed to %q", data), update.CallbackQuery.Message.Chat.ID)
		t.answerCallback(update)
	default:
		t.QuickSend(fmt.Sprintf("Unknown callback command %v", command), update.CallbackQuery.Message.Chat.ID)
	}
}

func (t Telegraphist) answerCallback(update tgbotapi.Update) {
	_, err := t.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err != nil {
		log.Println(err)
	}
}
