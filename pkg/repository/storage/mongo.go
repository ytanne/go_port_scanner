package storage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ytanne/go_nessus/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	client         *mongo.Client
	arpCollection  *mongo.Collection
	nmapCollection *mongo.Collection
	webCollection  *mongo.Collection
}

const (
	mongoConnectTimeout = time.Second * 10
)

func NewDatabaseRepository(cfg config.Config) (*mongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	mongoURI := fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongo db: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not ping mongo: %w", err)
	}

	clientDB := client.Database(cfg.Mongo.Database)

	newMongoDB := mongoDB{
		client:         client,
		arpCollection:  clientDB.Collection("arp_collection"),
		nmapCollection: clientDB.Collection("nmap_collection"),
		webCollection:  clientDB.Collection("web_collection"),
	}

	return &newMongoDB, nil
}

func (m *mongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("could not close client: %w", err)
	}

	return nil
}
