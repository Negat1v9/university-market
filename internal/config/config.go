package config

import (
	"os"
	"strconv"
	"strings"
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
	// admin telegram IDs
	AdminsIDs []int64
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
	adminIDsString := os.Getenv("ADMINS_IDS")
	if adminIDsString == "" {
		panic("no ADMINS_IDS")
	}

	return &WebCfg{
		CtxTimeOut: time.Second * time.Duration(stringToInt(ctxTimeOut)),
		Port:       port,
		JwtSecret:  []byte(jwtSecret),
		TgBotToken: tgToken,
		AdminsIDs:  stringToInt64Slice(adminIDsString),
	}
}

func stringToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic("error convert: " + s + " to int type")
	}
	return num
}

func stringToInt64Slice(s string) []int64 {
	slice := strings.Split(s, ",")
	res := make([]int64, 0, len(slice))
	for i := 0; i < len(slice); i++ {
		num, err := strconv.ParseInt(slice[i], 10, 64)
		if err != nil {
			panic("stringToInt64Slice: convert string to int64 s=" + slice[i])
		}
		res = append(res, num)
	}
	return res
}
