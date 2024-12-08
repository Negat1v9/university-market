package manager

import (
	"context"
	"log/slog"

	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	managerutils "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils"
	msgcrtr "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/msgcreater"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/static"
	tgbotmodel "github.com/Negat1v9/work-marketplace/model/tgBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *Manager) manageText(ctx context.Context, msg *tgbotapi.Message) {
	// skip update
	if m.cache.IsExist(msg.From.ID) {
		return
	}
	m.cache.Add(msg.From.ID)

	defer m.cache.Delete(msg.From.ID)

	cmd, err := m.store.TgCmd().Find(ctx, msg.From.ID)
	var response *tgbotapi.MessageConfig
	switch {
	case err == mongoStore.ErrNoTgCmd:
		// m.botClient.Send(msgcrtr.CreateTextMsg(msg.From.ID, static.NoTgCmdFinded))
		return
	case err != nil:
		m.log.Error("manage message", slog.String("err", err.Error()))
		m.botClient.Send(msgcrtr.CreateTextMsg(msg.From.ID, static.ErrBot))
		return
	}

	switch cmd.ExpectedAction {
	case tgbotmodel.WaitingForFiles:
		response, err = m.attachFilesTask(ctx, msg, cmd)
	default:
		return
	}
	if err != nil {
		m.log.Error("manage message", slog.String("err", err.Error()))
		m.botClient.Send(msgcrtr.CreateTextMsg(msg.From.ID, static.ErrBot))
	} else {
		m.botClient.Send(response)
	}
}

func (m *Manager) attachFilesTask(ctx context.Context, msg *tgbotapi.Message, tgCmd *tgbotmodel.UserCommand) (*tgbotapi.MessageConfig, error) {
	var fileID string
	switch {
	case msg.Document == nil:
		return msgcrtr.CreateTextMsg(msg.From.ID, static.FilesExpected), nil
	case msg.Document != nil:
		fileID = msg.Document.FileID
	}

	err := m.taskService.AttachFiles(ctx, tgCmd.TaskID, fileID)
	if err != nil {
		return nil, err
	}

	res := msgcrtr.CreateTextMsg(msg.From.ID, static.SuccessAttachFiles())
	res.ReplyMarkup = managerutils.CreateInlineOnAddFiles()

	return res, nil
}
