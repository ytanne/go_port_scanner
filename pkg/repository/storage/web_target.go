package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *mongoDB) CreateNewWebTarget(ctx context.Context, target entities.NmapTarget) (entities.NmapTarget, error) {
	if _, err := m.webCollection.InsertOne(ctx, target); err != nil {
		return entities.NmapTarget{}, fmt.Errorf("could not insert Nmap target. Error: %w.", err)
	}

	return target, fmt.Errorf("not implemented yet")
}

func (m *mongoDB) SaveWebResult(ctx context.Context, target entities.NmapTarget) (int, error) {
	if _, err := m.webCollection.UpdateOne(ctx, bson.M{"target": target.ARPscanID}, target); err != nil {
		return -1, fmt.Errorf("could not update result. Error: %w", err)
	}

	return -1, fmt.Errorf("not implemented yet")
}

func (m *mongoDB) RetrieveWebRecord(ctx context.Context, targetName string, id int) (entities.NmapTarget, error) {
	var result entities.NmapTarget

	if err := m.webCollection.FindOne(ctx, bson.M{"target": targetName}).Decode(&result); err != nil {
		return entities.NmapTarget{}, fmt.Errorf("could not find target %s. Error: %w", targetName, err)
	}

	return result, nil
}

func (m *mongoDB) RetrieveOldWebTargets(ctx context.Context, timelimit int) ([]entities.NmapTarget, error) {
	lastTime := time.Now().Add(time.Duration(-timelimit) * time.Minute)
	cursor, err := m.webCollection.Find(ctx, bson.M{"scanTime": bson.M{"lte": primitive.NewDateTimeFromTime(lastTime)}})
	if err != nil {
		return nil, fmt.Errorf("could not find targets with timelimit of %d. Error: %w", timelimit, err)
	}

	var results []entities.NmapTarget
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("could not get results from cursor. Error: %w", err)
	}

	return results, nil
}

func (m *mongoDB) RetrieveAllWebTargets(ctx context.Context) ([]entities.NmapTarget, error) {
	cursor, err := m.arpCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("could not find all targets. Error: %w", err)
	}

	var results []entities.NmapTarget
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("could not get results from cursor. Error: %w", err)
	}

	return results, nil
}
