package database

import (
	"FullStackOfYear/backend/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongoDB(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(config.Cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	MongoClient = client
	return client.Ping(ctx, nil)
}
