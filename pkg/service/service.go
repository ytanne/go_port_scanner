package service

import (
	"context"

	"github.com/ytanne/go_nessus/pkg/entities"
	"github.com/ytanne/go_nessus/pkg/repository"
)

type Communicate interface {
	SendMessage(msg string) error
	ReadMessage(ctx context.Context, msg chan string) error
}

type Store interface {
	CreateNewARPTarget(target string) (*entities.ARPTarget, error)
	SaveARPResult(target *entities.ARPTarget) (int, error)
	RetrieveARPRecord(target string) (*entities.ARPTarget, error)
	RetrieveOldARPTargets(timelimit int) ([]*entities.ARPTarget, error)
	RetrieveAllARPTargets() ([]*entities.ARPTarget, error)
	CreateNewNmapTarget(target string, id int) (*entities.NmapTarget, error)
	SaveNmapResult(target *entities.NmapTarget) (int, error)
	RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error)
	RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error)
	RetrieveAllNmapTargets() ([]*entities.NmapTarget, error)
	CreateNewWebTarget(target string, id int) (*entities.NmapTarget, error)
	SaveWebResult(target *entities.NmapTarget) (int, error)
	RetrieveWebRecord(target string, id int) (*entities.NmapTarget, error)
	RetrieveOldWebTargets(timelimit int) ([]*entities.NmapTarget, error)
	RetrieveAllWebTargets() ([]*entities.NmapTarget, error)
}

type NmapScan interface {
	ScanPorts(ctx context.Context, target string) ([]string, error)
	ScanWebPorts(ctx context.Context, target string) ([]string, error)
	ScanNetwork(ctx context.Context, target string) ([]string, error)
}

type Service struct {
	Communicate
	Store
	NmapScan
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Communicate: NewCommunicator(repo.Communicate),
		Store:       NewServiceStorage(repo.Store),
		NmapScan:    NewNmapScanner(repo.NmapScan),
	}
}
