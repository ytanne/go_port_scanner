package service

import (
	"context"

	"github.com/ytanne/go_port_scanner/pkg/entities"
)

type Keeper interface {
	CreateNewARPTarget(ctx context.Context, target entities.ARPTarget) (entities.ARPTarget, error)
	SaveARPResult(ctx context.Context, target entities.ARPTarget) (int, error)
	RetrieveARPRecord(ctx context.Context, target string) (entities.ARPTarget, error)
	RetrieveOldARPTargets(ctx context.Context, timelimit int) ([]entities.ARPTarget, error)
	RetrieveAllARPTargets(ctx context.Context) ([]entities.ARPTarget, error)

	CreateNewNmapTarget(ctx context.Context, target entities.NmapTarget) (entities.NmapTarget, error)
	SaveNmapResult(ctx context.Context, target entities.NmapTarget) (int, error)
	RetrieveNmapRecord(ctx context.Context, target string, id int) (entities.NmapTarget, error)
	RetrieveOldNmapTargets(ctx context.Context, timelimit int) ([]entities.NmapTarget, error)
	RetrieveAllNmapTargets(ctx context.Context) ([]entities.NmapTarget, error)

	CreateNewWebTarget(ctx context.Context, target entities.NmapTarget) (entities.NmapTarget, error)
	SaveWebResult(ctx context.Context, target entities.NmapTarget) (int, error)
	RetrieveWebRecord(ctx context.Context, target string, id int) (entities.NmapTarget, error)
	RetrieveOldWebTargets(ctx context.Context, timelimit int) ([]entities.NmapTarget, error)
	RetrieveAllWebTargets(ctx context.Context) ([]entities.NmapTarget, error)
}
