package service

import (
	"context"

	"github.com/ytanne/go_nessus/pkg/repository"
)

type CommunicateService struct {
	repo repository.Communicate
}

func NewCommunicator(repo repository.Communicate) *CommunicateService {
	return &CommunicateService{repo}
}

func (cm *CommunicateService) ReadMessage(ctx context.Context, msg chan string) error {
	return cm.repo.ReadMessage(ctx, msg)
}

func (cm *CommunicateService) SendMessage(msg string) error {
	return cm.repo.SendMessage(msg)
}
