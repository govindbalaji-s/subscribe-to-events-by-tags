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
	if _, ok := readUser(eventCreatorEmail, errorPrefix); !ok {
		return 1
	}
	insertResult, err := EventsCollection.InsertOne(Ctx, bson.D{
		{eventNameField, eventName},
		{eventVenueField, eventVenue},
		{eventTimeField, eventTime},
		{eventDurationField, eventDuration},
		/*{eventTagsField, eventTags},*/ // redundant since TagEvent also adds
		{eventSubscribersField, bson.A{eventCreatorEmail}},
		{eventCreatorField, eventCreatorEmail},
	})
	eventID := insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println(insertResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon insertion", err)
		return 2
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{emailField: eventCreatorEmail}, bson.D{
		{"$addToSet", bson.D{
			{userCreatedEventsField, eventID},
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
	if _, ok := readEvent(eventID, errorPrefix); !ok {
		return 1
	}
	updateResult, err := EventsCollection.UpdateOne(Ctx, bson.M{eventIDField: eventID}, eventDetailsMap)
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
	event, ok := readEvent(eventID, errorPrefix)
	if !ok {
		return 1
	}
	subscribers := event[eventSubscribersField].([]string)
	for _, userEmail := range subscribers {
		if SubscribeToEvent(eventID, userEmail, false) != 0 {
			return 2
		}
	}
	creatorEmail := event[eventCreatorField].(string)
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{emailField: creatorEmail}, bson.D{
		{"$pull", bson.D{
			{userCreatedEventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon updating createdEvents")
		return 3
	}
	tags := event[eventTagsField].([]string)
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
	if _, ok := readEvent(eventID, errorPrefix); !ok {
		return 1
	}
	if _, ok := readUser(userEmail, errorPrefix); !ok {
		return 2
	}
	updateResult, err := EventsCollection.UpdateOne(Ctx, bson.M{eventIDField: eventID}, bson.D{
		{updateOp, bson.D{
			{eventSubscribersField, userEmail},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, "upon updating EventsColln")
		return 3
	}
	updateResult, err = UsersCollection.UpdateOne(Ctx, bson.M{emailField: userEmail}, bson.D{
		{updateOp, bson.D{
			{userSubscribedEventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix, " upon updating UsersColln")
		return 4
	}
	return 0
}
