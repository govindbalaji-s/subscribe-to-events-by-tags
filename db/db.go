package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	Db     *mongo.Database
    Ctx     context.Context
)

const (
	uri    = "mongodb://127.0.0.1:27017"
	dbName = "set"
)

func Init() context.CancelFunc {
	Ctx := context.Background()
	Ctx, cancel := context.WithCancel(Ctx)
	Client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Errorf("set: couldn't connect to mongo: %v", err)
		return nil
	}
	err = Client.Connect(Ctx)
	if err != nil {
		fmt.Errorf("set: mongo client couldn't connect with background context: %v", err)
		return nil
	}
	Db = Client.Database(dbName)
	return cancel
}
