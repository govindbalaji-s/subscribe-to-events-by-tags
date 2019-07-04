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

	TimeFormat = "02-01-2006 15:04 (IST)"
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

// CreateTag creates a tag in the database with given name
// Return values:
// 0 - tag succesfully created
// 2 - tag already exists
// 3 - database error
func CreateTag(name string) int {
	errorPrefix := "set: db/tag.go: CreateTag:"
	if _, ok := ReadTag(name, errorPrefix); ok {
		return 2
	}
	insertResult, err := TagsCollection.InsertOne(Ctx, bson.D{
		{TagNameField, name},
		{FollowersField, bson.A{}},
		{EventsField, bson.A{}},
	})
	if err != nil {
		fmt.Println(errorPrefix, " name=", name, err)
		fmt.Println(insertResult)
		return 3
	}
	fmt.Println(insertResult)
	return 0
}

// TagEvent tags/untags the event of eventID with the tag of tagName. email is the email of logged in user
// toTag = true  => Tag, false => Untag
// Returns:
// 0 = successful
// 1 = tag does not exist
// 2 = event does not exist
// 3 = event is not created by the user tagging or not logged in
// 4 = db error on updating TagsCollection
// 5 = db error on updating EventsCollection
func TagEvent(email string, tagName string, eventID primitive.ObjectID, toTag bool) int {
	errorPrefix := "set: tag.go: TagEvent:"
	updateOp := "$addToSet"
	if !toTag {
		errorPrefix = "set: tag.go: TagEvent(untag):"
		updateOp = "$pull"
	}
	if _, ok := ReadTag(tagName, errorPrefix); !ok {
		return 1
	}

	event, ok := ReadEvent(eventID, errorPrefix)
	if !ok {
		return 2
	}

	if email != event[EventCreatorField].(string) {
		return 3
	}

	updateResult, err := TagsCollection.UpdateOne(Ctx, bson.M{TagNameField: tagName}, bson.D{
		{updateOp, bson.D{
			{EventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the tag with event", err)
		return 4
	}

	updateResult, err = EventsCollection.UpdateOne(Ctx, bson.M{EventIDField: eventID}, bson.D{
		{updateOp, bson.D{
			{TagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the event with tag", err)
		return 5
	}
	return 0
}

// FollowTag adds/removes the tag to list of followed tag and adds/removes the user to list of followers
// toFollow = true => follow, false => unfollow
// Returns:
// 0 - success
// 1 - user not found
// 2 - tag not found
// 3 - db update error in usersCollection
// 4 - db update error in tagsCollection
func FollowTag(email string, tagName string, toFollow bool) int {
	errorPrefix := "set: tag.go: FollowTag:"
	updateOp := "$addToSet"
	if !toFollow {
		errorPrefix = "set: tag.go: FollowTag(unfollow):"
		updateOp = "$pull"
	}
	if _, ok := ReadUser(email, errorPrefix); !ok {
		return 1
	}
	if _, ok := ReadTag(tagName, errorPrefix); !ok {
		return 2
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{EmailField: email}, bson.D{
		{updateOp, bson.D{
			{TagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating userCln", err)
		return 3
	}
	updateResult, err = TagsCollection.UpdateOne(Ctx, bson.M{TagNameField: tagName}, bson.D{
		{updateOp, bson.D{
			{FollowersField, email},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating tagCln", err)
		return 4
	}
	return 0
}

func SearchTags(query string) ([]bson.M, bool) {
	var tags []bson.M
	cur, err := TagsCollection.Find(Ctx, bson.D{
		{TagNameField, bson.D{
			{"$regex", query},
			{"$options", "i"},
		}},
	})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	defer cur.Close(Ctx)
	for cur.Next(Ctx) {
		tag := bson.M{}
		if err = cur.Decode(&tag); err != nil {
			fmt.Println(err)
			return nil, false
		}
		tags = append(tags, tag)
	}
	return tags, true
}

// CreateEvent creates an event in the db with the passed parameters.
// eventTime in TimeFormat
// Return values : code, eventID
// 0 = success
// 1 = creator DNE in usersColln
// 2 = db insert failed
// 3 = db update failed for user's created
// 4 = problem in updating tags
func CreateEvent(eventName string, eventVenue string, eventTime int64, eventDuration int64, eventTags []string, eventCreatorEmail string) (int, primitive.ObjectID) {
	errorPrefix := "set: event.go: CreateEvent:"
	if _, ok := ReadUser(eventCreatorEmail, errorPrefix); !ok {
		return 1, primitive.NilObjectID
	}
	insertResult, err := EventsCollection.InsertOne(Ctx, bson.D{
		{EventNameField, eventName},
		{EventVenueField, eventVenue},
		{EventTimeField, eventTime},
		{EventDurationField, eventDuration},
		/*{eventTagsField, eventTags},*/ // redundant since TagEvent also adds
		{EventSubscribersField, bson.A{eventCreatorEmail}},
		{EventCreatorField, eventCreatorEmail},
	})
	eventID := insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println(insertResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon insertion", err)
		return 2, primitive.NilObjectID
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{EmailField: eventCreatorEmail}, bson.D{
		{"$addToSet", bson.D{
			{UserCreatedEventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon updating UsersColln")
		return 3, primitive.NilObjectID
	}
	//to deal with tags
	for _, tagName := range eventTags {
		if TagEvent(eventCreatorEmail, tagName, eventID, true) != 0 {
			return 4, primitive.NilObjectID
		}
	}
	return 0, eventID
}

// EditEventDetails edits the event with given eventid and sets the fields present in the eventDetailsMap
// fields assumed to be a subset of {name, venue, time, duration}
// Return values :
// 0 = success
// 1 = event dne
// 2 = creator is not the user
// 3 = db update failed
func EditEventDetails(eventID primitive.ObjectID, eventDetailsMap bson.M, userEmail string) int {
	errorPrefix := "set: event.go: EditEventDetails:"
	event, ok := ReadEvent(eventID, errorPrefix)
	if !ok {
		return 1
	}
	if event[EventCreatorField].(string) != userEmail {
		return 2
	}
	updateResult, err := EventsCollection.UpdateOne(Ctx, bson.M{EventIDField: eventID}, bson.D{
		{"$set", eventDetailsMap},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon update")
		return 3
	}
	return 0
}

//DeleteEvent deletes the event, after clearing up tags and subscribers and creator
// Returns:
// 0 = success
// 1 = event dne
// 2 = creator not signed in
// 3 = error in clearing up subs
// 4 = error in clearing up the creator
// 5 = error in clearing up the tags
// 6 = error in deleting event
func DeleteEvent(eventID primitive.ObjectID, userEmail string) int {
	errorPrefix := "set: event.go: DeleteEvent:"
	event, ok := ReadEvent(eventID, errorPrefix)
	if !ok {
		return 1
	}
	if event[EventCreatorField].(string) != userEmail {
		return 2
	}
	fmt.Println(event)
	subscribers := event[EventSubscribersField].(primitive.A)
	for _, userEmail := range subscribers {
		if SubscribeToEvent(eventID, userEmail.(string), false) != 0 {
			return 3
		}
	}
	creatorEmail := event[EventCreatorField].(string)
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{EmailField: creatorEmail}, bson.D{
		{"$pull", bson.D{
			{UserCreatedEventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon updating createdEvents")
		return 4
	}
	tags := event[EventTagsField].(primitive.A)
	for _, tagName := range tags {
		if TagEvent(creatorEmail, tagName.(string), eventID, false) != 0 {
			return 5
		}
	}
	_, err = EventsCollection.DeleteOne(Ctx, bson.M{EventIDField: eventID})
	if err != nil {
		fmt.Println(errorPrefix, "upon deleting event")
		return 6
	}
	return 0
}

//SubscribeToEvent makes the user sub/unsub to the event.
// toSubsribe = true => subscribe else unsubscribe
// Return values:
// 0 = success
// 1 = event dne
// 2 = user dne
// 3 = db error in EventsColln
// 4 = db error in UsersColln
func SubscribeToEvent(eventID primitive.ObjectID, userEmail string, toSubscribe bool) int {
	errorPrefix := "set: event.go: SubscribeToEvent:"
	updateOp := "$addToSet"
	if !toSubscribe {
		errorPrefix = "set: event.go: SubscribeToEvent(unsub):"
		updateOp = "$pull"
	}
	if _, ok := ReadEvent(eventID, errorPrefix); !ok {
		return 1
	}
	if _, ok := ReadUser(userEmail, errorPrefix); !ok {
		return 2
	}
	updateResult, err := EventsCollection.UpdateOne(Ctx, bson.M{EventIDField: eventID}, bson.D{
		{updateOp, bson.D{
			{EventSubscribersField, userEmail},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon updating EventsColln")
		return 3
	}
	updateResult, err = UsersCollection.UpdateOne(Ctx, bson.M{EmailField: userEmail}, bson.D{
		{updateOp, bson.D{
			{UserSubscribedEventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, " upon updating UsersColln")
		return 4
	}
	return 0
}

func SearchEventsByName(query string) ([]bson.M, bool) {
	var events []bson.M
	cur, err := EventsCollection.Find(Ctx, bson.D{
		{EventNameField, bson.D{
			{"$regex", query},
			{"$options", "i"},
		}},
	})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	defer cur.Close(Ctx)
	for cur.Next(Ctx) {
		event := bson.M{}
		if err = cur.Decode(&event); err != nil {
			fmt.Println(err)
			return nil, false
		}
		events = append(events, event)
	}
	return events, true
}
