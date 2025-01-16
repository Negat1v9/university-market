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

	var response tgbotapi.Chattable
	var err error

	switch {
	case managerutils.PostTaskCallBack == cb.Data:
		response, err = m.publichTask(ctx, cb)
	case managerutils.ShareContactCallBack == cb.Data:
		response, err = m.shareContact(ctx, cb)
	case managerutils.CreateEventCallBack == cb.Data:
		response, err = m.onCreateEvent(ctx, cb)
	case managerutils.StartSendingMessages == cb.Data:
		response, err = m.startSendingMessages(ctx, cb)
	}

	if err != nil {
		m.log.Error("manage callback", slog.String("err", err.Error()))
		m.botClient.Send(msgcrtr.CreateTextMsg(cb.From.ID, static.ErrBot))
	} else {
		m.botClient.Send(response)
	}

}

func (m *Manager) publichTask(ctx context.Context, cb *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	cmd, err := m.store.TgCmd().Find(ctx, cb.From.ID)
	switch {
	case err == mongoStore.ErrNoTgCmd:
		return msgcrtr.CreateTextMsg(cb.From.ID, static.NoTgCmdFinded), nil
	case err != nil:
		return msgcrtr.CreateTextMsg(cb.From.ID, static.ErrBot), nil
	}

	err = m.taskService.PublishTask(ctx, cmd.TaskID)
	if err != nil {
		return nil, err
	}
	res := msgcrtr.CreateTextMsg(cb.From.ID, static.SuccessPublichTask)
	res.ReplyMarkup = managerutils.CreateInlineOnPublichTask(m.webAppBaseUrl, cmd.TaskID)
	return res, nil
}

func (m *Manager) shareContact(ctx context.Context, cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	cmd, err := m.store.TgCmd().Find(ctx, cb.From.ID)
	switch {
	case err == mongoStore.ErrNoTgCmd:
		return msgcrtr.CreateTextMsg(cb.From.ID, static.NoTgCmdFinded), nil
	case err != nil:
		return msgcrtr.CreateTextMsg(cb.From.ID, static.ErrBot), nil
	}

	if cb.From.UserName == "" {
		return msgcrtr.CreateTextMsg(cb.From.ID, static.ErrNoUserName), nil
	}
	userInfo, err := m.userService.User(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	msgForUser := msgcrtr.CreateTextMsg(userInfo.TelegramID, static.MsgShareContact(cb.From.UserName))
	msgForUser.ReplyMarkup = managerutils.CreateInlineOnPublichTask(m.webAppBaseUrl, cmd.TaskID)
	err = m.botClient.Send(msgForUser)
	if err != nil {
		return nil, err
	}
	markup := managerutils.CreateInlineAfterShareContact(m.webAppBaseUrl, cmd.TaskID)

	return tgbotapi.NewEditMessageReplyMarkup(cb.From.ID, cb.Message.MessageID, *markup), nil
}

func (m *Manager) onCreateEvent(ctx context.Context, cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	err := m.store.TgCmd().Delete(ctx, cb.From.ID)
	if err != nil && err != mongoStore.ErrNoTgCmd {
		m.log.Error("Manager.onCreateEvent.TgCmd.Delete", slog.String("err", err.Error()))
		return nil, err
	}
	tgCmd := &tgbotmodel.UserCommand{
		ID:             cb.From.ID,
		ExpectedAction: tgbotmodel.WaitingEventCaption,
	}

	err = m.store.TgCmd().Create(ctx, tgCmd)
	if err != nil && err != mongoStore.ErrNoTgCmd {
		m.log.Error("Manager.onCreateEvent.TgCmd.Create", slog.String("err", err.Error()))
		return nil, err
	}
	msg := msgcrtr.CreateTextMsg(cb.From.ID, static.CreateEvent)

	return msg, nil
}

func (m *Manager) startSendingMessages(ctx context.Context, cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	cmd, err := m.store.TgCmd().Find(ctx, cb.From.ID)
	switch {
	case err == mongoStore.ErrNoTgCmd:
		return msgcrtr.CreateTextMsg(cb.From.ID, static.NoTgCmdFinded), nil
	case err != nil:
		return msgcrtr.CreateTextMsg(cb.From.ID, static.ErrBot), nil
	}

	if err = m.adminService.StartSendingEvent(ctx, cmd.EventID); err != nil {
		return nil, err
	}

	return msgcrtr.CreateTextMsg(cb.From.ID, static.OnStartSendEvent), nil
}
