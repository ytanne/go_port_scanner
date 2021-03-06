package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ytanne/go_port_scanner/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *mongoDB) CreateNewARPTarget(ctx context.Context, target entities.ARPTarget) (entities.ARPTarget, error) {
	if _, err := m.arpCollection.InsertOne(ctx, target); err != nil {
		return entities.ARPTarget{}, fmt.Errorf("could not insert ARP target. Error: %w", err)
	}

	return target, nil
}

func (m *mongoDB) SaveARPResult(ctx context.Context, target entities.ARPTarget) (int, error) {
	update := bson.M{
		"$set": target,
	}

	if _, err := m.arpCollection.UpdateOne(ctx, bson.M{"target": target.Target}, update); err != nil {
		return -1, fmt.Errorf("could not update result. Error: %w", err)
	}

	return target.ID, nil
}

func (m *mongoDB) RetrieveARPRecord(ctx context.Context, targetName string) (entities.ARPTarget, error) {
	var result entities.ARPTarget

	if err := m.arpCollection.FindOne(ctx, bson.M{"target": targetName}).Decode(&result); err != nil {
		return entities.ARPTarget{}, fmt.Errorf("could not find target %s. Error: %w", targetName, err)
	}

	return result, nil
}

func (m *mongoDB) RetrieveOldARPTargets(ctx context.Context, timelimit int) ([]entities.ARPTarget, error) {
	lastTime := time.Now().Add(time.Duration(-timelimit) * time.Minute)
	cursor, err := m.arpCollection.Find(ctx, bson.M{"scan_time": bson.M{"$lte": lastTime}})
	if err != nil {
		return nil, fmt.Errorf("could not find targets with timelimit of %d. Error: %w", timelimit, err)
	}

	var results []entities.ARPTarget
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("could not get results from cursor. Error: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no old arp targets found")
	}

	return results, nil
}

func (m *mongoDB) RetrieveAllARPTargets(ctx context.Context) ([]entities.ARPTarget, error) {
	cursor, err := m.arpCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("could not find all targets. Error: %w", err)
	}

	var results []entities.ARPTarget
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("could not get results from cursor. Error: %w", err)
	}

	return results, nil
}
