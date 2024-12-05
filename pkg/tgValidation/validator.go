package tgvalidation

import (
	"encoding/json"
	"fmt"
	"net/url"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var (
	contactKeyValue = "contact"
)

type ReqPhoneResponse struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	UserID      int64  `json:"user_id"`
}

// package for validating data that came from telegram mini app

func ValidateInitData(initData string, tgBotToken string) error {
	// expIn := 24 * time.Hour
	return initdata.Validate(initData, tgBotToken, -1)
}

func ParseInitData(initData string) (initdata.InitData, error) {
	return initdata.Parse(initData)
}

func ParsePhoneNumber(response string) (ReqPhoneResponse, error) {
	q, err := url.ParseQuery(response)
	if err != nil {
		return ReqPhoneResponse{}, fmt.Errorf("invalid response string")
	}
	data := q.Get(contactKeyValue)
	if data == "" {
		return ReqPhoneResponse{}, fmt.Errorf("invalid response string")
	}
	var r ReqPhoneResponse
	err = json.Unmarshal([]byte(data), &r)
	if err != nil {
		return ReqPhoneResponse{}, err
	}
	if r.PhoneNumber == "" {
		return ReqPhoneResponse{}, fmt.Errorf("no phone number")
	}

	return r, nil
}
