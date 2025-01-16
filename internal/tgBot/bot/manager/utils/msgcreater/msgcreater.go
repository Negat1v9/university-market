package msgcrtr

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateTextMsg(id int64, text string) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "html"
	msg.DisableWebPagePreview = true
	return &msg
}

func CreatePhotoMsg(id int64, fileID, caption string) *tgbotapi.PhotoConfig {
	msg := tgbotapi.NewPhoto(id, tgbotapi.FileID(fileID))
	msg.ParseMode = "html"
	msg.Caption = caption
	return &msg
}
