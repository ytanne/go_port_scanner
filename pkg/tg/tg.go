package tg

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	NoUpdates = "No updates from bot"
)

type Telegram struct {
	Bot       *tgbotapi.BotAPI
	TGconfigs *tgbotapi.UpdateConfig
	GroupID   int64
}

func NewTelegramConn(APItoken string, GroupID int64) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(APItoken)
	if err != nil {
		return nil, err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &Telegram{
		Bot:       bot,
		TGconfigs: &u,
		GroupID:   GroupID,
	}, nil
}

func (t *Telegram) SendMessage(msg string) error {
	reply := tgbotapi.NewMessage(t.GroupID, msg)

	_, err := t.Bot.Send(reply)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

func (t *Telegram) ReadMessages(ctx context.Context, msg chan string) error {
	updates, err := t.Bot.GetUpdatesChan(*t.TGconfigs)
	if err != nil {
		return fmt.Errorf("could not get updates from channel %d: %w", t.GroupID, err)
	}

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				msg <- update.Message.Text
			}
		case <-ctx.Done():
			t.Bot.StopReceivingUpdates()
			close(msg)
			return nil
		case <-time.After(time.Second * 10):
			// case for no updates
		}
	}
}
