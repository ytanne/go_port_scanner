package storage

import (
	"context"
	"math/rand"
	"time"

	"github.com/ytanne/go_port_scanner/pkg/entities"
	"github.com/ytanne/go_port_scanner/pkg/service"
)

type serviceStorage struct {
	repo service.Keeper
}

func NewDatabaseService(repo service.Keeper) *serviceStorage {
	return &serviceStorage{
		repo: repo,
	}
}

func (ss *serviceStorage) CreateNewARPTarget(ctx context.Context, target entities.ARPTarget) (entities.ARPTarget, error) {
	target.ID = rand.Int()
	target.ScanTime = time.Now()

	return ss.repo.CreateNewARPTarget(ctx, target)
}

func (ss *serviceStorage) RetrieveARPRecord(ctx context.Context, target string) (entities.ARPTarget, error) {
	return ss.repo.RetrieveARPRecord(ctx, target)
}

func (ss *serviceStorage) SaveARPResult(ctx context.Context, target entities.ARPTarget) (int, error) {
	target.ScanTime = time.Now()
	target.ScanTime = time.Now()

	return ss.repo.SaveARPResult(ctx, target)
}

func (ss *serviceStorage) RetrieveOldARPTargets(ctx context.Context, timelimit int) ([]entities.ARPTarget, error) {
	return ss.repo.RetrieveOldARPTargets(ctx, timelimit)
}

func (ss *serviceStorage) RetrieveAllARPTargets(ctx context.Context) ([]entities.ARPTarget, error) {
	return ss.repo.RetrieveAllARPTargets(ctx)
}

func (ss *serviceStorage) CreateNewNmapTarget(ctx context.Context, target entities.NmapTarget, id int) (entities.NmapTarget, error) {
	target.ARPscanID = id
	target.ID = rand.Int()
	target.ScanTime = time.Now()

	return ss.repo.CreateNewNmapTarget(ctx, target)
}

func (ss *serviceStorage) RetrieveNmapRecord(ctx context.Context, target string, id int) (entities.NmapTarget, error) {
	return ss.repo.RetrieveNmapRecord(ctx, target, id)
}

func (ss *serviceStorage) SaveNmapResult(ctx context.Context, target entities.NmapTarget) (int, error) {
	target.ScanTime = time.Now()
	target.ScanTime = time.Now()

	return ss.repo.SaveNmapResult(ctx, target)
}

func (ss *serviceStorage) RetrieveOldNmapTargets(ctx context.Context, timelimit int) ([]entities.NmapTarget, error) {
	return ss.repo.RetrieveOldNmapTargets(ctx, timelimit)
}

func (ss *serviceStorage) RetrieveAllNmapTargets(ctx context.Context) ([]entities.NmapTarget, error) {
	return ss.repo.RetrieveAllNmapTargets(ctx)
}

func (ss *serviceStorage) CreateNewWebTarget(ctx context.Context, target entities.NmapTarget, id int) (entities.NmapTarget, error) {
	target.ARPscanID = id
	target.ID = rand.Int()
	target.ScanTime = time.Now()

	return ss.repo.CreateNewWebTarget(ctx, target)
}

func (ss *serviceStorage) RetrieveWebRecord(ctx context.Context, target string, id int) (entities.NmapTarget, error) {
	return ss.repo.RetrieveWebRecord(ctx, target, id)
}

func (ss *serviceStorage) SaveWebResult(ctx context.Context, target entities.NmapTarget) (int, error) {
	target.ScanTime = time.Now()

	return ss.repo.SaveWebResult(ctx, target)
}

func (ss *serviceStorage) RetrieveOldWebTargets(ctx context.Context, timelimit int) ([]entities.NmapTarget, error) {
	return ss.repo.RetrieveOldWebTargets(ctx, timelimit)
}

func (ss *serviceStorage) RetrieveAllWebTargets(ctx context.Context) ([]entities.NmapTarget, error) {
	return ss.repo.RetrieveAllWebTargets(ctx)
}
