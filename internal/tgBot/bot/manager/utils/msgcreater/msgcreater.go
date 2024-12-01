package msgcrtr

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateTextMsg(id int64, text string) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "html"
	msg.DisableWebPagePreview = true
	return &msg
}
