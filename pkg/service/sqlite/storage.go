package sqlite

import (
	"github.com/ytanne/go_nessus/pkg/entities"
	"github.com/ytanne/go_nessus/pkg/repository/sqlite"
)

type NMAP interface {
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

type DBKeeper interface {
	NMAP
}

type serviceStorage struct {
	repo sqlite.DBKeeper
}

func NewDatabaseService(repo sqlite.DBKeeper) DBKeeper {
	return &serviceStorage{
		repo: repo,
	}
}

func (ss *serviceStorage) CreateNewARPTarget(target string) (*entities.ARPTarget, error) {
	return ss.repo.CreateNewARPTarget(target)
}

func (ss *serviceStorage) RetrieveARPRecord(target string) (*entities.ARPTarget, error) {
	return ss.repo.RetrieveARPRecord(target)
}

func (ss *serviceStorage) SaveARPResult(target *entities.ARPTarget) (int, error) {
	return ss.repo.SaveARPResult(target)
}

func (ss *serviceStorage) RetrieveOldARPTargets(timelimit int) ([]*entities.ARPTarget, error) {
	return ss.repo.RetrieveOldARPTargets(timelimit)
}

func (ss *serviceStorage) RetrieveAllARPTargets() ([]*entities.ARPTarget, error) {
	return ss.repo.RetrieveAllARPTargets()
}

func (ss *serviceStorage) CreateNewNmapTarget(target string, id int) (*entities.NmapTarget, error) {
	return ss.repo.CreateNewNmapTarget(target, id)
}

func (ss *serviceStorage) RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error) {
	return ss.repo.RetrieveNmapRecord(target, id)
}

func (ss *serviceStorage) SaveNmapResult(target *entities.NmapTarget) (int, error) {
	return ss.repo.SaveNmapResult(target)
}

func (ss *serviceStorage) RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error) {
	return ss.repo.RetrieveOldNmapTargets(timelimit)
}

func (ss *serviceStorage) RetrieveAllNmapTargets() ([]*entities.NmapTarget, error) {
	return ss.repo.RetrieveAllNmapTargets()
}

func (ss *serviceStorage) CreateNewWebTarget(target string, id int) (*entities.NmapTarget, error) {
	return ss.repo.CreateNewWebTarget(target, id)
}

func (ss *serviceStorage) RetrieveWebRecord(target string, id int) (*entities.NmapTarget, error) {
	return ss.repo.RetrieveWebRecord(target, id)
}

func (ss *serviceStorage) SaveWebResult(target *entities.NmapTarget) (int, error) {
	return ss.repo.SaveWebResult(target)
}

func (ss *serviceStorage) RetrieveOldWebTargets(timelimit int) ([]*entities.NmapTarget, error) {
	return ss.repo.RetrieveOldWebTargets(timelimit)
}

func (ss *serviceStorage) RetrieveAllWebTargets() ([]*entities.NmapTarget, error) {
	return ss.repo.RetrieveAllWebTargets()
}
