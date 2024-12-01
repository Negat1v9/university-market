package config

import (
	"time"
)

func NewConfigMock() *Config {

	return &Config{
		Env: "MOCK",
		MongoConfing: &MongoCfg{
			MongoUrl:    "URL",
			MongoDbName: "Test Name",
		},
		BotConfig: &BotCfg{
			BotToken:      "token",
			WebAppBaseUrl: "http://none",
			NumberWorkers: 10,
		},
		WebConfig: &WebCfg{
			CtxTimeOut: 10 * time.Second,
			Port:       ":8080",
			JwtSecret:  []byte("Secret"),
			TgBotToken: "token",
		},
	}
}
