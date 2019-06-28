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
	Ctx    context.Context

	UsersCollection,
	TagsCollection,
	EventsCollection *mongo.Collection
)

const (
	uri            = "mongodb://127.0.0.1:27017"
	dbName         = "set"
	emailField     = "email"
	eventIDField   = "_id"
	tagNameField   = "name"
	tagsField      = "tags"
	eventsField    = "events"
	followersField = "followers"
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
	UsersCollection = Db.Collection("users")
	TagsCollection = Db.Collection("tags")
	EventsCollection = Db.Collection("events")
	return cancel
}
