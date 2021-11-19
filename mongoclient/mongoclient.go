package mongoclient

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoHandler struct {
	Client   *mongo.Client
	Database string
}

func NewHandler(address, database string) (*MongoHandler, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(address))
	if err != nil {
		return nil, err
	}
	return &MongoHandler{
		Client:   cl,
		Database: database,
	}, err
}
