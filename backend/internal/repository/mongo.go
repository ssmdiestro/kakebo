package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

func NewMongo(ctx context.Context, uri string) (*Mongo, error) {
	if uri == "" {
		return nil, errors.New("MONGO_URI vacío")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo.Ping: %w", err)
	}
	return &Mongo{Client: client}, nil
}

func (m *Mongo) Collection(db, coll string) *mongo.Collection {
	return m.Client.Database(db).Collection(coll)
}
