package telegram

import (
	"log"

	_ "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Communicator interface {
	SendMessage(msg string) error
	ReadMessage(msg chan string) error
}

type tgCommunicate struct {
	Bot       *tgbotapi.BotAPI
	TGconfigs *tgbotapi.UpdateConfig
	GroupID   int64
}

func NewCommunicatorRepository(APItoken string, gi int64) (Communicator, error) {
	log.Printf("Obtained API token: %s", APItoken)
	log.Printf("Group ID: %d", gi)
	bot, err := tgbotapi.NewBotAPI(APItoken)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return &tgCommunicate{bot, &u, gi}, nil
}

func (r *tgCommunicate) ReadMessage(msg chan string) error {
	return r.readMessages(msg)
}

func (r *tgCommunicate) SendMessage(msg string) error {
	return r.sendMessage(msg)
}

func (t *tgCommunicate) sendMessage(msg string) error {
	reply := tgbotapi.NewMessage(t.GroupID, msg)
	_, err := t.Bot.Send(reply)
	if err != nil {
		_, err = t.Bot.Send(reply)
	}
	return err
}

func (t *tgCommunicate) readMessages(msg chan string) error {
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
