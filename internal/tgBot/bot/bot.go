package bot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Negat1v9/work-marketplace/internal/config"
	"github.com/Negat1v9/work-marketplace/internal/services"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	cfg       *config.BotCfg
	log       *slog.Logger
	manager   *manager.Manager
	botClient tgbot.TgBotClient
}

func New(cfg *config.BotCfg, log *slog.Logger, c tgbot.TgBotClient, services *services.Services, store storage.Store) *Bot {
	return &Bot{
		cfg:       cfg,
		log:       log,
		manager:   manager.New(log, c, cfg.WebAppBaseUrl, services, store),
		botClient: c,
	}
}

func (b *Bot) Start() {
	updates := b.botClient.UpdatesChan()
	b.skipUpdateOnStart(updates)
	for i := 0; i <= b.cfg.NumberWorkers; i++ {
		go func() {
			for {
				select {
				case update := <-updates:
					b.manager.UpdateBot(update)
				default:
					time.Sleep(time.Millisecond * 200)
				}
			}
		}()
	}

	b.log.Info("bot workers", slog.Int("count", b.cfg.NumberWorkers))
}

// skips updates that could come while the application is not working
func (b *Bot) skipUpdateOnStart(updates tgbotapi.UpdatesChannel) {
	b.log.Info("bot start skip updates")
	timer := time.NewTimer(time.Second * 10)

	for {
		select {
		case <-timer.C:
			b.log.Info("bot stop skip updates")
			return
		case upd := <-updates:
			if upd.Message != nil {
				if upd.Message.SuccessfulPayment != nil {
					b.log.Warn(fmt.Sprintf(
						"successful payment from userID %d in the amount is %d of tg stars",
						upd.Message.From.ID, upd.Message.SuccessfulPayment.TotalAmount,
					))
				}
			}
		}
	}
}
