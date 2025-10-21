package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Insert[T any](list []T, tableName string, databaseName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	table := MongoClient.Database(databaseName).Collection(tableName)
	var resultList []any
	for _, v := range list {
		var i interface{}
		reqBodyBytes := new(bytes.Buffer)
		err := json.NewEncoder(reqBodyBytes).Encode(v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(reqBodyBytes.Bytes(), &i)
		if err != nil {
			return err
		}
		resultList = append(resultList, i)
	}
	if len(resultList) > 0 {
		_, err := table.InsertMany(ctx, resultList)
		if err != nil {
			fmt.Println("Insert - table.InsertMany(): ", err)
			return err
		}
	}
	return nil
}

func SelectWithError[T any](resultList []T, criteria bson.M, tableName string, databaseName string) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	table := MongoClient.Database(databaseName).Collection(tableName)
	cursor, err := table.Find(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("Select - table.Find(): %w", err)
	}

	for cursor.Next(ctx) {
		var object T
		err := cursor.Decode(&object)
		if err != nil {
			return nil, fmt.Errorf("Select - cursor.Decode(): %w", err)
		}
		resultList = append(resultList, object)
	}

	return resultList, nil
}

func SetMongoConnection(uri string) error {
	if uri == "" {
		return errors.New("MONGO_URI vacío")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("mongo.Connect: %w", err)
	}

	// Verificación real de conexión
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongo.Ping: %w", err)
	}

	MongoClient = client
	fmt.Println("MongoDB connection established")
	return nil
}
