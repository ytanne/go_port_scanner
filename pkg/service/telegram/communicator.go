package telegram

import (
	"context"

	"github.com/ytanne/go_nessus/pkg/repository/telegram"
)

type Communicator interface {
	SendMessage(msg string) error
	ReadMessage(ctx context.Context, msg chan string) error
}

type communicateService struct {
	repo telegram.Communicator
}

func NewCommunicatorService(repo telegram.Communicator) Communicator {
	return &communicateService{repo}
}

func (cm *communicateService) ReadMessage(ctx context.Context, msg chan string) error {
	return cm.repo.ReadMessage(ctx, msg)
}

func (cm *communicateService) SendMessage(msg string) error {
	return cm.repo.SendMessage(msg)
}
