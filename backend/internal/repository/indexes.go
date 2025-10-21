package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureRecordIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: map[string]int{"date": 1}},
		{Keys: map[string]int{"year": 1, "contableMonth": 1, "weekOfMonth": 1}},
		{Keys: map[string]int{"subcategory.category": 1}},
	}
	_, err := coll.Indexes().CreateMany(ctx, models, &options.CreateIndexesOptions{})
	if err != nil {
		return fmt.Errorf("CreateIndexes: %w", err)
	}
	return nil
}
