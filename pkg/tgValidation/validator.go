package tgvalidation

import (
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// package for validating data that came from telegram mini app

func ValidateInitData(initData string, tgBotToken string) error {
	// expIn := 24 * time.Hour
	return initdata.Validate(initData, tgBotToken, -1)
}

func ParseInitData(initData string) (initdata.InitData, error) {
	return initdata.Parse(initData)
}
