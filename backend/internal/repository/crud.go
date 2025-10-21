package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne
func InsertOne[T any](ctx context.Context, coll *mongo.Collection, doc T, opts ...*options.InsertOneOptions) (interface{}, error) {
	res, err := coll.InsertOne(ctx, doc, opts...)
	if err != nil {
		return nil, fmt.Errorf("InsertOne: %w", err)
	}
	return res.InsertedID, nil
}

// InsertMany
func InsertMany[T any](ctx context.Context, coll *mongo.Collection, docs []T, opts ...*options.InsertManyOptions) ([]interface{}, error) {
	if len(docs) == 0 {
		return nil, nil
	}
	bulk := make([]interface{}, len(docs))
	for i := range docs {
		bulk[i] = docs[i]
	}
	res, err := coll.InsertMany(ctx, bulk, opts...)
	if err != nil {
		return nil, fmt.Errorf("InsertMany: %w", err)
	}
	return res.InsertedIDs, nil
}

// Find (slice)
func Find[T any](ctx context.Context, coll *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	cur, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("Find: %w", err)
	}
	defer cur.Close(ctx)

	var out []T
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("Cursor.All: %w", err)
	}
	return out, nil
}

// FindOne (un solo doc)
func FindOne[T any](ctx context.Context, coll *mongo.Collection, filter interface{}, opts ...*options.FindOneOptions) (*T, error) {
	var out T
	if err := coll.FindOne(ctx, filter, opts...).Decode(&out); err != nil {
		return nil, fmt.Errorf("FindOne: %w", err)
	}
	return &out, nil
}

// UpdateOne
func UpdateOne(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	res, err := coll.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return 0, fmt.Errorf("UpdateOne: %w", err)
	}
	return res.ModifiedCount, nil
}

// ReplaceOne
func ReplaceOne[T any](ctx context.Context, coll *mongo.Collection, filter interface{}, doc T, opts ...*options.ReplaceOptions) (int64, error) {
	res, err := coll.ReplaceOne(ctx, filter, doc, opts...)
	if err != nil {
		return 0, fmt.Errorf("ReplaceOne: %w", err)
	}
	return res.ModifiedCount, nil
}

// DeleteOne
func DeleteOne(ctx context.Context, coll *mongo.Collection, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	res, err := coll.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, fmt.Errorf("DeleteOne: %w", err)
	}
	return res.DeletedCount, nil
}
