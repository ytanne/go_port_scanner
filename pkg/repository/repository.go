package repository

import (
	"database/sql"

	"github.com/ytanne/go_nessus/pkg/entities"
	"github.com/ytanne/go_nessus/pkg/tg"
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
	RetrieveAllARPTargets() ([]*entities.ARPTarget, error)
	CreateNewNmapTarget(target string, id int) (*entities.NmapTarget, error)
	SaveNmapResult(target *entities.NmapTarget) (int, error)
	RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error)
	RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error)
	RetrieveAllNmapTargets() ([]*entities.NmapTarget, error)
}

type NmapScan interface {
	ScanPorts(target string) ([]byte, error)
	ScanNetwork(target string) ([]byte, error)
}

type Nessus interface {
	ListScans() (*entities.ScanList, error)
}

type Repository struct {
	Store
	Communicate
	NmapScan
	Nessus
}

func NewRepository(db *sql.DB, t *tg.Telegram, AccessKey, SecretKey, URL string) *Repository {
	return &Repository{
		Store:       NewDatabase(db),
		Communicate: NewCommunicator(t),
		NmapScan:    NewScanner(),
		Nessus:      NewNessusClient(AccessKey, SecretKey, URL),
	}
}
