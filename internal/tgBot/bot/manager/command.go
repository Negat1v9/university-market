package manager

import (
	"context"
	"log/slog"
	"strings"

	msgcrtr "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/msgcreater"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/static"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	startCmd = "/start"
	helpCmd  = "/help"
)

// returns true if the text is a command (starts with /)
func isCommand(text string) bool {
	if len(text) > 0 {
		return text[0] == '/'
	}

	return false
}

func (m *Manager) manageCommand(ctx context.Context, msg *tgbotapi.Message) {
	var res *tgbotapi.MessageConfig
	var err error
	switch {
	case strings.HasPrefix(msg.Text, startCmd):
		res, err = m.isStartCmd(ctx, msg)
	case msg.Text == helpCmd:
		res = m.isHelpCmd(msg.From.ID)
	// do nothing on unknown cmd
	default:
		return
	}

	if err != nil {
		res = msgcrtr.CreateTextMsg(msg.From.ID, static.ErrBot)
	}
	if err = m.botClient.Send(res); err != nil {
		m.log.Error("bot send cmd", slog.String("err", err.Error()))
	}

}

func (m *Manager) isStartCmd(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	referralID := referallIDFromStartCmd(msg.Text)
	newUser := usermodel.NewUser(
		msg.From.ID,
		msg.From.UserName,
		msg.From.FirstName+" "+msg.From.LastName,
		referralID,
	)
	_, err := m.userService.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return msgcrtr.CreateTextMsg(msg.From.ID, static.MsgStartCommand(msg.From.UserName)), nil
}

func (m *Manager) isHelpCmd(userID int64) *tgbotapi.MessageConfig {
	return msgcrtr.CreateTextMsg(userID, static.HelpCmd)
}

func referallIDFromStartCmd(s string) int64 {
	startParam := strings.TrimPrefix(s, startCmd+" ")
	referralID, err := utils.ConvertStringToInt64(startParam)
	if err != nil {
		return 0
	}
	return referralID
}
