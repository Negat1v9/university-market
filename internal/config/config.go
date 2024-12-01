package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env          string
	MongoConfing *MongoCfg
	BotConfig    *BotCfg
	WebConfig    *WebCfg
}

func NewConfig() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		panic("no ENV")
	}
	webCfg := newWebCfg()
	mongoCfg := newMongoCfg()
	botCfg := newBotCfg()

	return &Config{
		Env:          env,
		MongoConfing: mongoCfg,
		BotConfig:    botCfg,
		WebConfig:    webCfg,
	}
}

type MongoCfg struct {
	MongoUrl    string
	MongoDbName string
}

func newMongoCfg() *MongoCfg {
	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		panic("no MONGO_URL")
	}
	mongoDbName := os.Getenv("MONGO_DB_NAME")
	if mongoUrl == "" {
		panic("no MONGO_DB_NAME")
	}
	return &MongoCfg{
		MongoUrl:    mongoUrl,
		MongoDbName: mongoDbName,
	}
}

type BotCfg struct {
	BotToken      string
	WebAppBaseUrl string
	NumberWorkers int
}

func newBotCfg() *BotCfg {
	botToken := os.Getenv("TG_BOT_TOKEN")
	if botToken == "" {
		panic("no TG_BOT_TOKEN")
	}
	numberWorkers := os.Getenv("NUMBER_WORKER")
	if numberWorkers == "" {
		panic("no NUMBER_WORKER")
	}
	webAppBaseUrl := os.Getenv("WEB_APP_BASE_URL")
	if webAppBaseUrl == "" {
		panic("no WEB_APP_BASE_URL")
	}
	return &BotCfg{
		BotToken:      botToken,
		WebAppBaseUrl: webAppBaseUrl,
		NumberWorkers: stringToInt(numberWorkers),
	}
}

type WebCfg struct {
	// base context time out for requests
	CtxTimeOut time.Duration
	// Server listen port
	Port string
	// jwt secret key for generate tokens
	JwtSecret []byte
	// telegram bot token from botFather
	TgBotToken string
}

func newWebCfg() *WebCfg {
	ctxTimeOut := os.Getenv("CONTEXT_TIMEOUT")
	if ctxTimeOut == "" {
		panic("no CONTEXT_TIMEOUT")
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		panic("no SERVER_PORT")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("no JWT_SECRET")
	}
	tgToken := os.Getenv("TG_BOT_TOKEN")
	if tgToken == "" {
		panic("no TG_BOT_TOKEN")
	}
	return &WebCfg{
		CtxTimeOut: time.Second * time.Duration(stringToInt(ctxTimeOut)),
		Port:       port,
		JwtSecret:  []byte(jwtSecret),
		TgBotToken: tgToken,
	}
}

func stringToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic("error convert: " + s + " to int type")
	}
	return num
}
