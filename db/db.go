package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Uri                       = "mongodb://127.0.0.1:27017"
	DbName                    = "set"
	EmailField                = "email"
	EventIDField              = "_id"
	TagNameField              = "name"
	TagsField                 = "tags"
	EventsField               = "events"
	FollowersField            = "followers"
	EventNameField            = "name"
	EventVenueField           = "venue"
	EventTimeField            = "time"
	EventDurationField        = "duration"
	EventTagsField            = "tags"
	EventSubscribersField     = "subscribers"
	EventCreatorField         = "creator"
	UserCreatedEventsField    = "createdEvents"
	UserSubscribedEventsField = "subscribedEvents"
)

func Init() context.CancelFunc {
	Ctx := context.Background()
	Ctx, cancel := context.WithCancel(Ctx)
	Client, err := mongo.NewClient(options.Client().ApplyURI(Uri))
	if err != nil {
		fmt.Errorf("set: couldn't connect to mongo: %v", err)
		return nil
	}
	err = Client.Connect(Ctx)
	if err != nil {
		fmt.Errorf("set: mongo client couldn't connect with background context: %v", err)
		return nil
	}
	Db = Client.Database(DbName)
	UsersCollection = Db.Collection("users")
	TagsCollection = Db.Collection("tags")
	EventsCollection = Db.Collection("events")
	return cancel
}

func ReadTag(tagName string, errorPrefix string) (tag bson.M, ok bool) {
	tag = bson.M{}
	err := TagsCollection.FindOne(Ctx, bson.M{
		TagNameField: tagName,
	}).Decode(&tag)
	if err != nil {
		fmt.Println(errorPrefix, "tagName =", tagName, err)
		return nil, false
	}
	return tag, true
}

func ReadEvent(eventID primitive.ObjectID, errorPrefix string) (event bson.M, ok bool) {
	event = bson.M{}
	err := EventsCollection.FindOne(Ctx, bson.M{
		EventIDField: eventID,
	}).Decode(&event)
	if err != nil {
		fmt.Println(errorPrefix, "eventID =", eventID, err)
		return nil, false
	}
	return event, true
}

func ReadUser(email string, errorPrefix string) (user bson.M, ok bool) {
	user = bson.M{}
	err := UsersCollection.FindOne(Ctx, bson.M{
		EmailField: email,
	}).Decode(&user)
	if err != nil {
		fmt.Println(errorPrefix, "email =", email, err)
		return nil, false
	}
	return user, true
}
