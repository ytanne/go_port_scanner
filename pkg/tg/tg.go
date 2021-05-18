package tg

import (
	"log"

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
	log.Printf("Obtained API token: %s", APItoken)
	log.Printf("Group ID: %d", GroupID)
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
		_, err = t.Bot.Send(reply)
	}
	return err
}

func (t *Telegram) ReadMessages(msg chan string) error {
	updates, err := t.Bot.GetUpdatesChan(*t.TGconfigs)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message != nil {
			msg <- update.Message.Text //fmt.Sprintf("%s: %s", update.Message.From, update.Message.Text)
		}
	}

	close(msg)
	return nil
}
