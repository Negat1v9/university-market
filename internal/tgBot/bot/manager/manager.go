package manager

import (
	"context"
	"log/slog"
	"time"

	"github.com/Negat1v9/work-marketplace/internal/services"
	paymentservice "github.com/Negat1v9/work-marketplace/internal/services/payment"
	taskservice "github.com/Negat1v9/work-marketplace/internal/services/task"
	userservice "github.com/Negat1v9/work-marketplace/internal/services/user"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/cache"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	fiveSecondTimeOut = time.Second * 5
)

type Manager struct {
	webAppBaseUrl  string
	log            *slog.Logger
	cache          *cache.BotCache
	paymentService paymentservice.PaymentService
	taskService    taskservice.TaskService
	userService    userservice.UserService
	botClient      tgbot.TgBotClient
	store          storage.Store
}

func New(log *slog.Logger, bc tgbot.TgBotClient, webAppUrl string, services *services.Services, store storage.Store) *Manager {
	return &Manager{
		webAppBaseUrl:  webAppUrl,
		log:            log,
		cache:          cache.New(),
		paymentService: services.PaymentService,
		taskService:    services.TaskService,
		userService:    services.UserService,
		botClient:      bc,
		store:          store,
	}
}

func (m *Manager) UpdateBot(update tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondTimeOut)
	defer cancel()

	switch {
	case update.PreCheckoutQuery != nil:
		m.precheckOut(ctx, update.PreCheckoutQuery)

	case update.Message != nil:
		switch {
		// check on content successfulPayment type in message
		case update.Message.SuccessfulPayment != nil:
			m.successPayment(ctx, update.Message)
		case isCommand(update.Message.Text):
			m.manageCommand(ctx, update.Message)

		default:
			m.manageText(ctx, update.Message)
		}

	case update.CallbackQuery != nil:
		m.manageCallBack(ctx, update.CallbackQuery)
	}
}
