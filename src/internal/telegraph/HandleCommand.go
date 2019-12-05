package telegraph

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
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
			msg.Text = "List of commands:\n" +
				"/screen - takes a screenshot"
			_, err := t.bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
		case "screen":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Imagine that that's a screenshot")
			_, err := t.bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
		default:
			_, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command"))
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
