package composites

import (
	"github.com/ytanne/go_nessus/pkg/app"
	"github.com/ytanne/go_nessus/pkg/config"
	repo "github.com/ytanne/go_nessus/pkg/repository/storage"
	serv "github.com/ytanne/go_nessus/pkg/service/storage"
)

type MongoComposite struct {
	Serv app.Keeper
}

func NewMongoComposite(cfg config.Config) (MongoComposite, error) {
	dbRepo, err := repo.NewDatabaseRepository(cfg)
	if err != nil {
		return MongoComposite{}, err
	}

	servDB := serv.NewDatabaseService(dbRepo)

	return MongoComposite{servDB}, nil
}
