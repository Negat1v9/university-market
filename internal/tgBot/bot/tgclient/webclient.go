package tgclient

import (
	"context"

	managerutils "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils"
	msgcrtr "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/msgcreater"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/static"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Client) SendRespond(ctx context.Context, tgCreaterID int64, workerID string, task *taskmodel.Task) error {
	text := static.OnRespondFromWorker()
	if task.Meta != nil {
		text += static.AddInformationTask(task.Meta)
	}
	response := msgcrtr.CreateTextMsg(tgCreaterID, text)
	response.ReplyMarkup = managerutils.CreateInlineRespondOnTask(c.cfg.WebAppBaseUrl, task.ID, workerID)
	// response.ReplyMarkup = managerutils.CreateInlineOnPublichTask(c.cfg.WebAppBaseUrl, tas)

	return c.Send(response)
}

func (c *Client) WaitFiles(ctx context.Context, tgCreaterID int64) error {
	response := msgcrtr.CreateTextMsg(tgCreaterID, static.WaitingFiles)
	response.ReplyMarkup = managerutils.CreateInlineOnAddFiles()

	return c.Send(response)
}

func (c *Client) SelectWorker(ctx context.Context, tgWorkerID int64, task *taskmodel.Task) error {
	text := static.WorkerSelected
	if task.Meta != nil {
		text += static.AddInformationTask(task.Meta)
	}
	response := msgcrtr.CreateTextMsg(tgWorkerID, text)
	response.ReplyMarkup = managerutils.CreateInlineOnWorkerSelected(c.cfg.WebAppBaseUrl, task.ID)

	return c.Send(response)
}

func (c *Client) SendFiles(ctx context.Context, tgWorkerID int64, files []string) error {

	groupFiles := make([]interface{}, 0, len(files))
	for i := 0; i < len(files); i++ {
		media := tgbotapi.NewInputMediaDocument(tgbotapi.FileID(files[i]))
		groupFiles = append(groupFiles, media)
	}
	msg := tgbotapi.NewMediaGroup(tgWorkerID, groupFiles)

	return c.SendGroupMedia(msg)
}
