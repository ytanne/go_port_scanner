package composites

import (
	"fmt"

	"github.com/ytanne/go_port_scanner/pkg/app"
	"github.com/ytanne/go_port_scanner/pkg/config"
	serv "github.com/ytanne/go_port_scanner/pkg/service/communication"
)

type CommunicationComposite struct {
	Serv app.Communicator
}

func NewCommunicationComposite(cfg config.Config) (CommunicationComposite, error) {
	commServ, err := serv.NewDiscordBot(cfg.Discord.Token)
	if err != nil {
		return CommunicationComposite{}, fmt.Errorf("Could not create new database repository: %w", err)
	}

	return CommunicationComposite{
		Serv: commServ,
	}, nil
}
