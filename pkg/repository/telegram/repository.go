package telegram

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Communicator interface {
	SendMessage(msg string) error
	ReadMessage(ctx context.Context, msg chan string) error
}

type tgCommunicate struct {
	Bot       *tgbotapi.BotAPI
	TGconfigs *tgbotapi.UpdateConfig
	GroupID   int64
}

func NewCommunicatorRepository(APItoken string, gi int64) (Communicator, error) {
	bot, err := tgbotapi.NewBotAPI(APItoken)
	if err != nil {
		return nil, err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return &tgCommunicate{bot, &u, gi}, nil
}

func (r *tgCommunicate) ReadMessage(ctx context.Context, msg chan string) error {
	return r.readMessages(ctx, msg)
}

func (r *tgCommunicate) SendMessage(msg string) error {
	return r.sendMessage(msg)
}

func (t *tgCommunicate) sendMessage(msg string) error {
	reply := tgbotapi.NewMessage(t.GroupID, msg)

	_, err := t.Bot.Send(reply)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

func (t *tgCommunicate) readMessages(ctx context.Context, msg chan string) error {
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
