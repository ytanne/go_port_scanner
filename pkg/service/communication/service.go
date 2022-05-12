package communication

import (
	"github.com/ytanne/go_nessus/pkg/models"
	"github.com/ytanne/go_nessus/pkg/service"
)

type communicationService struct {
	repo service.Communicator
}

func NewCommunicationService(discordRepo service.Communicator) communicationService {
	return communicationService{
		repo: discordRepo,
	}
}

func (c communicationService) SendMessage(msg, channelID string) error {
	return c.repo.SendMessage(msg, channelID)
}

func (c communicationService) MessageReadChannel() chan models.Message {
	return c.repo.MessageReadChannel()
}
