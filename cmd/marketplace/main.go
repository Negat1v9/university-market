package main

import (
	"context"
	"log"

	"github.com/Negat1v9/work-marketplace/internal/app"
	"github.com/Negat1v9/work-marketplace/internal/config"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	"github.com/Negat1v9/work-marketplace/pkg/logger"
)

func main() {
	cfg := config.NewConfig()
	log.Println("config is loaded")

	logger := logger.New(cfg.Env)

	db, err := mongoStore.NewClient(context.Background(), cfg.MongoConfing.MongoUrl, cfg.MongoConfing.MongoDbName)
	if err != nil {
		panic(err.Error())
	}

	store := mongoStore.New(db)

	app := app.New(cfg, logger, store)

	if err = app.RunApp(); err != nil {
		logger.Warn("app is down")
	}
}
