package repository

import (
	"context"

	"github.com/ytanne/go_nessus/pkg/tg"
)

type Communicator struct {
	t *tg.Telegram
}

func NewCommunicator(t *tg.Telegram) *Communicator {
	return &Communicator{t}
}

func (c *Communicator) ReadMessage(ctx context.Context, msg chan string) error {
	return c.t.ReadMessages(ctx, msg)
}

func (c *Communicator) SendMessage(msg string) error {
	return c.t.SendMessage(msg)
}
