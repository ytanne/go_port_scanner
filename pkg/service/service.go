package service

import (
	"github.com/ytanne/go_nessus/pkg/entities"
	"github.com/ytanne/go_nessus/pkg/repository"
)

type Communicate interface {
	SendMessage(msg string) error
	ReadMessage(msg chan string) error
}

type Store interface {
	CreateNewARPTarget(target string) (*entities.ARPTarget, error)
	SaveARPResult(target *entities.ARPTarget) (int, error)
	RetrieveARPRecord(target string) (*entities.ARPTarget, error)
	RetrieveOldARPTargets(timelimit int) ([]*entities.ARPTarget, error)
	CreateNewNmapTarget(target string, id int) (*entities.NmapTarget, error)
	SaveNmapResult(target *entities.NmapTarget) (int, error)
	RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error)
	RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error)
}

type NmapScan interface {
	ScanPorts(target string) ([]string, error)
	ScanNetwork(target string) ([]string, error)
}

type ARPScan interface {
}

type Nessus interface {
	ListScans() (*entities.ScanList, error)
}

type Service struct {
	Communicate
	Store
	NmapScan
	Nessus
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Communicate: NewCommunicator(repo.Communicate),
		Store:       NewServiceStorage(repo.Store),
		NmapScan:    NewNmapScanner(repo.NmapScan),
		Nessus:      NewNessusService(repo.Nessus),
	}
}
