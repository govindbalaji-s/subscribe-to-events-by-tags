package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEvent creates an event in the db with the passed parameters.
// Return values :
// 0 = success
// 1 = creator DNE in usersColln
// 2 = db insert failed
// 3 = db update failed for user's created
// 4 = problem in updating tags
func CreateEvent(eventName string, eventVenue string, eventTime time.Time, eventDuration time.Duration, eventTags []string, eventCreatorEmail string) int {
	errorPrefix := "set: event.go: CreateEvent:"
	if _, ok := ReadUser(eventCreatorEmail, errorPrefix); !ok {
		return 1
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
		return 2
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{EmailField: eventCreatorEmail}, bson.D{
		{"$addToSet", bson.D{
			{UserCreatedEventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon updating UsersColln")
		return 3
	}
	//to deal with tags
	for _, tagName := range eventTags {
		if TagEvent(eventCreatorEmail, tagName, eventID, true) != 0 {
			return 4
		}
	}
	return 0
}

// EditEventDetails edits the event with given eventid and sets the fields present in the eventDetailsMap
// fields assumed to be a subset of {name, venue, time, duration}
// Return values :
// 0 = success
// 1 = event dne
// 2 = db update failed
func EditEventDetails(eventID primitive.ObjectID, eventDetailsMap bson.M) int {
	errorPrefix := "set: event.go: EditEventDetails:"
	if _, ok := ReadEvent(eventID, errorPrefix); !ok {
		return 1
	}
	updateResult, err := EventsCollection.UpdateOne(Ctx, bson.M{EventIDField: eventID}, eventDetailsMap)
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon update")
		return 2
	}
	return 0
}

//DeleteEvent deletes the event, after clearing up tags and subscribers and creator
// Returns:
// 0 = success
// 1 = event dne
// 2 = error in clearing up subs
// 3 = error in clearing up the creator
// 4 = error in clearing up the tags
func DeleteEvent(eventID primitive.ObjectID) int {
	errorPrefix := "set: event.go: DeleteEvent:"
	event, ok := ReadEvent(eventID, errorPrefix)
	if !ok {
		return 1
	}
	subscribers := event[EventSubscribersField].([]string)
	for _, userEmail := range subscribers {
		if SubscribeToEvent(eventID, userEmail, false) != 0 {
			return 2
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
		return 3
	}
	tags := event[EventTagsField].([]string)
	for _, tagName := range tags {
		if TagEvent(creatorEmail, tagName, eventID, false) != 0 {
			return 4
		}
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
