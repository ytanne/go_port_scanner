package repository

import (
	"github.com/ytanne/go_nessus/pkg/tg"
)

type Communicator struct {
	t *tg.Telegram
}

func NewCommunicator(t *tg.Telegram) *Communicator {
	return &Communicator{t}
}

func (c *Communicator) ReadMessage(msg chan string) error {
	return c.t.ReadMessages(msg)
}

func (c *Communicator) SendMessage(msg string) error {
	return c.t.SendMessage(msg)
}
