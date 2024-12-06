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

func (m *Manager) manageCallBack(ctx context.Context, cb *tgbotapi.CallbackQuery) {
	// skip update
	if m.cache.IsExist(cb.From.ID) {
		return
	}
	m.cache.Add(cb.From.ID)
	defer m.cache.Delete(cb.From.ID)

	cmd, err := m.store.TgCmd().Find(ctx, cb.From.ID)
	var response tgbotapi.Chattable
	switch {
	case err == mongoStore.ErrNoTgCmd:
		m.botClient.Send(msgcrtr.CreateTextMsg(cb.From.ID, static.NoTgCmdFinded))
		return
	case err != nil:
		m.botClient.Send(msgcrtr.CreateTextMsg(cb.From.ID, static.ErrBot))
		return
	}

	switch {
	case managerutils.PostTaskCallBack == cb.Data && cmd.ExpectedAction == tgbotmodel.WaitingForFiles:
		response, err = m.publichTask(ctx, cmd.TaskID, cb)
	case managerutils.ShareContactCallBack == cb.Data && cmd.ExpectedAction == tgbotmodel.WorkerShareContact:
		response, err = m.shareContact(ctx, cmd, cb)
	}

	if err != nil {
		m.log.Error("manage callback", slog.String("err", err.Error()))
		m.botClient.Send(msgcrtr.CreateTextMsg(cb.From.ID, static.ErrBot))
	} else {
		m.botClient.Send(response)
	}

}

func (m *Manager) publichTask(ctx context.Context, taskID string, cb *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	err := m.taskService.PublishTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	res := msgcrtr.CreateTextMsg(cb.From.ID, static.SuccessPublichTask)
	res.ReplyMarkup = managerutils.CreateInlineOnPublichTask(m.webAppBaseUrl, taskID)
	return res, nil
}

func (m *Manager) shareContact(ctx context.Context, cmd *tgbotmodel.UserCommand, cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	if cb.From.UserName == "" {
		return msgcrtr.CreateTextMsg(cb.From.ID, static.ErrNoUserName), nil
	}
	userInfo, err := m.userService.User(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	msgForUser := msgcrtr.CreateTextMsg(userInfo.TelegramID, static.MsgShareContact(cb.From.UserName))

	err = m.botClient.Send(msgForUser)
	if err != nil {
		return nil, err
	}
	markup := managerutils.CreateInlineAfterShareContact(m.webAppBaseUrl, cmd.TaskID)

	return tgbotapi.NewEditMessageReplyMarkup(cb.From.ID, cb.Message.MessageID, *markup), nil
}
