package bot

import (
	"log/slog"
	"time"

	"github.com/Negat1v9/work-marketplace/internal/config"
	"github.com/Negat1v9/work-marketplace/internal/services"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/manager"
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
