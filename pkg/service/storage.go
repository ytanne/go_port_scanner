package service

import (
	"github.com/ytanne/go_nessus/pkg/entities"
	"github.com/ytanne/go_nessus/pkg/repository"
)

type serviceStorage struct {
	repo repository.Store
}

func NewServiceStorage(repo repository.Store) *serviceStorage {
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
