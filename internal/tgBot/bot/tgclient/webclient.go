package tgclient

import (
	"context"

	managerutils "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils"
	msgcrtr "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/msgcreater"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/static"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

func (c *Client) SendRespond(ctx context.Context, tgCreaterID int64, taskID string, worker *usermodel.User) error {

	response := msgcrtr.CreateTextMsg(tgCreaterID, static.OnRespondFromWorker(worker.WorkerInfo))
	response.ReplyMarkup = managerutils.CreateInlineRespondOnTask(c.cfg.WebAppBaseUrl, taskID, worker.ID)

	return c.Send(response)
}

func (c *Client) WaitFiles(ctx context.Context, tgCreaterID int64) error {
	response := msgcrtr.CreateTextMsg(tgCreaterID, static.WaitingFiles)
	response.ReplyMarkup = managerutils.CreateInlineOnAddFiles()

	return c.Send(response)
}

func (c *Client) SelectWorker(ctx context.Context, tgWorkerID int64, taskID string) error {
	response := msgcrtr.CreateTextMsg(tgWorkerID, static.WorkerSelected)
	response.ReplyMarkup = managerutils.CreateInlineOnWorkerSelected(c.cfg.WebAppBaseUrl, taskID)

	return c.Send(response)
}

// func (c *Client) SendFiles(ctx context.Context, tgWorkerID int64, files []string) error {

// 	tgbotapi.NewMediaGroup(tgWorkerID)
// 	return nil
// }
