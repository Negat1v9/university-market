package tgbot

import (
	"context"

	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotClient interface {
	Send(tgbotapi.Chattable) error
	Request(*tgbotapi.PreCheckoutConfig) error
	UpdatesChan() tgbotapi.UpdatesChannel
}

type WebTgClient interface {
	SendRespond(ctx context.Context, tgCreaterID int64, workerID string, task *taskmodel.Task) error
	WaitFiles(ctx context.Context, tgCreaterID int64) error
	SelectWorker(ctx context.Context, tgWorkerID int64, task *taskmodel.Task) error
	SendFiles(ctx context.Context, tgWorkerID int64, files []string) error
	SendEventMsg(userTgID int64, event *eventmodel.Event) error
}
