package app

import (
	"log/slog"

	"github.com/Negat1v9/work-marketplace/internal/config"
	"github.com/Negat1v9/work-marketplace/internal/services"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot"
	"github.com/Negat1v9/work-marketplace/internal/tgBot/bot/tgclient"
	"github.com/Negat1v9/work-marketplace/internal/web"
)

type App struct {
	cfg   *config.Config
	log   *slog.Logger
	store storage.Store
}

func New(cfg *config.Config, log *slog.Logger, store storage.Store) *App {
	return &App{
		cfg:   cfg,
		log:   log,
		store: store,
	}
}
func (a *App) RunApp() error {
	botClient := tgclient.NewBot(a.cfg.BotConfig)

	services := services.NewServiceBuilder(a.cfg, a.log, botClient, a.store)

	server := web.NewServer(a.cfg, a.log, services, a.store)

	bot := bot.New(a.cfg.BotConfig, a.log, botClient, services, a.store)

	bot.Start()
	a.log.Info("start bot")

	a.log.Info("start server")
	return server.Run()

}
