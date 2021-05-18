package service

import (
	"github.com/ytanne/go_nessus/pkg/entities"
	"github.com/ytanne/go_nessus/pkg/repository"
)

type NessusService struct {
	repo repository.Nessus
}

func NewNessusService(repo repository.Nessus) *NessusService {
	return &NessusService{repo}
}

func (ns *NessusService) ListScans() (*entities.ScanList, error) {
	return nil, nil
}
