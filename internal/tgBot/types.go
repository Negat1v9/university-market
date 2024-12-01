package tgbot

import (
	"context"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotClient interface {
	Send(tgbotapi.Chattable) error
	Request(*tgbotapi.PreCheckoutConfig) error
	UpdatesChan() tgbotapi.UpdatesChannel
}

type WebTgClient interface {
	SendRespond(ctx context.Context, tgCreaterID int64, taskID string, worker *usermodel.User) error
	WaitFiles(ctx context.Context, tgCreaterID int64) error
	SelectWorker(ctx context.Context, tgWorkerID int64, taskID string) error
}
