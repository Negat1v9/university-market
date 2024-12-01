package manager

import (
	"context"
	"log/slog"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	msgcrtr "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/msgcreater"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager/utils/static"
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *Manager) precheckOut(ctx context.Context, msg *tgbotapi.PreCheckoutQuery) {
	response := tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: msg.ID,
	}
	payment, err := m.store.Payment().Find(ctx, filters.New().Add(filters.PaymentByID(msg.InvoicePayload)).Filters())
	if err != nil {
		m.log.Error("precheck out", slog.String("err", err.Error()))
		response.OK = false
		response.ErrorMessage = static.ErrorPreCheckOutAmount
		err = m.botClient.Request(&response)
		if err != nil {
			m.log.Error("send precheck out", slog.String("err", err.Error()))
			return
		}
	}
	if payment.Amount != msg.TotalAmount {
		response.OK = false
		response.ErrorMessage = static.ErrorPreCheckOutAmount
		err = m.botClient.Request(&response)
		if err != nil {
			m.log.Error("send precheck out", slog.String("err", err.Error()))
		}
		return
	}
	response.OK = true
	err = m.botClient.Request(&response)
	if err != nil {
		m.log.Error("send precheck out", slog.String("err", err.Error()))
	}
}

func (m *Manager) successPayment(ctx context.Context, msg *tgbotapi.Message) {
	sp := paymentmodel.SuccessPayment{
		PaymentID:   msg.SuccessfulPayment.InvoicePayload,
		Amount:      msg.SuccessfulPayment.TotalAmount,
		TgPaymentID: msg.SuccessfulPayment.TelegramPaymentChargeID,
	}
	var res *tgbotapi.MessageConfig
	err := m.paymentService.SuccessBalancePayment(ctx, &sp)
	if err != nil {
		m.log.Error("success payment", slog.String("payment id", msg.SuccessfulPayment.TelegramPaymentChargeID))
		res = msgcrtr.CreateTextMsg(msg.From.ID, static.MsgErrOnSuccessPayment(msg.SuccessfulPayment.TelegramPaymentChargeID))
	} else {
		res = msgcrtr.CreateTextMsg(msg.From.ID, static.BalancePayment(msg.SuccessfulPayment.TotalAmount))
	}

	err = m.botClient.Send(res)
	if err != nil {
		m.log.Error("success payment", slog.String("err", err.Error()))
	}
}
