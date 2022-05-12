package composites

import (
	"fmt"

	"github.com/ytanne/go_nessus/pkg/app"
	"github.com/ytanne/go_nessus/pkg/config"
	repo "github.com/ytanne/go_nessus/pkg/repository/storage"
	serv "github.com/ytanne/go_nessus/pkg/service/storage"
)

type DBComposite struct {
	DBServ app.Keeper
}

func NewDBComposite(cfg config.Config) (DBComposite, error) {
	dbRepo, err := repo.NewDatabaseRepository(cfg)
	if err != nil {
		return DBComposite{}, fmt.Errorf("Could not create new database repository: %w", err)
	}

	dbServ := serv.NewDatabaseService(dbRepo)

	return DBComposite{
		DBServ: dbServ,
	}, nil
}
