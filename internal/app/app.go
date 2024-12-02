package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil {
			a.log.Error("app is down", slog.String("err", err.Error()))
		}
	}()

	<-sigs
	now := time.Now().UTC()
	a.log.Info("get signal stop app", slog.String("time UTC", now.Format(time.DateTime)))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		a.log.Error("server shotdown", slog.String("err", err.Error()))
	}

	<-ctx.Done()

	a.log.Info("stop app", slog.Duration("shotdown time", time.Duration(time.Now().UTC().Sub(now).Seconds())))

	return nil
}
