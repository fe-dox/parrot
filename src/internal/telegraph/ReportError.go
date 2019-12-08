package telegraph

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (t Telegraphist) ReportError(message string, chatID int64) {
	t.bot.Send(tgbotapi.NewMessage(chatID, message))
}
