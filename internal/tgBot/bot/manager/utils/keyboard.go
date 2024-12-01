package managerutils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	PostTaskCallBack     = "1"
	ShareContactCallBack = "2"
)

func CreateInlineOnAddFiles() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Опубликовать", PostTaskCallBack),
		),
	)
	return &kb
}

func CreateInlineRespondOnTask(baseUrl, taskID, workerID string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("ℹ️ профиль работника", tgbotapi.WebAppInfo{
				URL: fmt.Sprintf("%s/worker/info?taskID=%s&workerID=%s", baseUrl, taskID, workerID),
			}),
		),
	)
	return &kb
}

func CreateInlineOnPublichTask(baseUrl, taskID string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("Работа", tgbotapi.WebAppInfo{
				URL: fmt.Sprintf("%s/task?taskID=%s", baseUrl, taskID),
			}),
		),
	)
	return &kb
}

func CreateInlineOnWorkerSelected(baseUrl, taskID string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Дать контакт пользоватю", ShareContactCallBack),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("Работа", tgbotapi.WebAppInfo{
				URL: fmt.Sprintf("%s/worker/task?taskID=%s", baseUrl, taskID),
			}),
		),
	)
	return &kb
}
