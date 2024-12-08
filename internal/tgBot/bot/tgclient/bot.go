package tgclient

import (
	"time"

	"github.com/Negat1v9/work-marketplace/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	buffer = 100
)

const (
	timeOutBot         = 5
	skipLastUpdateTime = time.Second * 5
)

type Client struct {
	cfg    *config.BotCfg
	client *tgbotapi.BotAPI
}

// Info: NewBot return tg bot client if cfg contain invalid token panic
func NewBot(cfg *config.BotCfg) *Client {
	client, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		panic(err.Error())
	}
	client.Buffer = buffer

	b := &Client{
		cfg:    cfg,
		client: client,
	}

	return b
}

func (b *Client) Send(res tgbotapi.Chattable) error {

	_, err := b.client.Send(res)

	return err
}

// Info: use only for answer precheck out query
func (b *Client) Request(preCheckOut *tgbotapi.PreCheckoutConfig) error {
	_, err := b.client.Request(preCheckOut)
	return err
}

func (b *Client) SendGroupMedia(msg tgbotapi.MediaGroupConfig) error {
	_, err := b.client.SendMediaGroup(msg)
	return err
}

func (b *Client) UpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)

	u.Timeout = timeOutBot

	updates := b.client.GetUpdatesChan(u)

	return updates
}
