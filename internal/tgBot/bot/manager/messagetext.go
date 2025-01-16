package manager

import (
	"context"
	"log/slog"
	"time"

	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	managerutils "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils"
	msgcrtr "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/msgcreater"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/static"
	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
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
	var response tgbotapi.Chattable
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
	case tgbotmodel.WaitingEventCaption:
		response, err = m.createEvent(ctx, msg)
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

func (m *Manager) createEvent(ctx context.Context, msg *tgbotapi.Message) (tgbotapi.Chattable, error) {
	var event eventmodel.Event
	var res tgbotapi.Chattable
	switch {
	case len(msg.Photo) != 0:
		event = eventmodel.Event{
			WithImage: true,
			CreatorID: msg.From.ID,
			UserType:  "all",
			Caption:   msg.Caption,
			FileID:    msg.Photo[0].FileID,
			CreatedAt: time.Now().UTC(),
		}
		msg := msgcrtr.CreatePhotoMsg(msg.From.ID, msg.Photo[0].FileID, msg.Caption)
		msg.ReplyMarkup = managerutils.CreateInlineStartMailingList()
		res = msg
	case msg.Text != "":
		event = eventmodel.Event{
			WithImage: false,
			CreatorID: msg.From.ID,
			UserType:  "all",
			Caption:   msg.Text,
			CreatedAt: time.Now().UTC(),
		}
		msg := msgcrtr.CreateTextMsg(msg.From.ID, msg.Text)
		msg.ReplyMarkup = managerutils.CreateInlineStartMailingList()
		res = msg
	default:
		return msgcrtr.CreateTextMsg(msg.From.ID, static.ErrTypeOnCreateEvent), nil
	}

	eventID, err := m.adminService.CreateEvent(ctx, &event)
	if err != nil {
		return nil, err
	}

	if err = m.store.TgCmd().Delete(ctx, msg.From.ID); err != nil {
		m.log.Error("Manager.createEvent.TgCmd.Delete", slog.String("err", err.Error()))
		return nil, err
	}

	cmd := tgbotmodel.UserCommand{
		ID:             msg.From.ID,
		ExpectedAction: tgbotmodel.WaitingStartSendEvent,
		EventID:        eventID,
	}
	if err = m.store.TgCmd().Create(ctx, &cmd); err != nil {
		m.log.Error("Manager.createEvent.TgCmd.Create", slog.String("err", err.Error()))
		return nil, err
	}

	return res, nil
}
