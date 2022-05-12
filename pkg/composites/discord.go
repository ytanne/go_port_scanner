package composites

import (
	"fmt"

	"github.com/ytanne/go_nessus/pkg/app"
	"github.com/ytanne/go_nessus/pkg/config"
	repo "github.com/ytanne/go_nessus/pkg/repository/communication"
	serv "github.com/ytanne/go_nessus/pkg/service/communication"
)

type CommunicationComposite struct {
	Serv app.Communicator
}

func NewCommunicationComposite(cfg config.Config) (CommunicationComposite, error) {
	communicationRepo, err := repo.NewDiscordBot(cfg.Discord.Token)
	if err != nil {
		return CommunicationComposite{}, fmt.Errorf("Could not create new database repository: %w", err)
	}

	discordServ := serv.NewCommunicationService(communicationRepo)

	return CommunicationComposite{
		Serv: discordServ,
	}, nil
}
