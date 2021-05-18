package service

import (
	"github.com/ytanne/go_nessus/pkg/repository"
)

type CommunicateService struct {
	repo repository.Communicate
}

func NewCommunicator(repo repository.Communicate) *CommunicateService {
	return &CommunicateService{repo}
}

func (cm *CommunicateService) ReadMessage(msg chan string) error {
	return cm.repo.ReadMessage(msg)
}

func (cm *CommunicateService) SendMessage(msg string) error {
	return cm.repo.SendMessage(msg)
}
