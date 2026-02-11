package db

import (
	"context"
	"crud/internal/config"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, nil, fmt.Errorf("Mongo connection failed: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("Mongo ping failed: %w", err)
	}

	database := client.Database(cfg.MongoDbName)

	return client, database, nil

}

func DisConnect(client *mongo.Client) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	return client.Disconnect(ctx)
}
