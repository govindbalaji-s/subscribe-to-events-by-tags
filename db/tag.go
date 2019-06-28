package db

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTag creates a tag in the database with given name
// Return values:
// 0 - tag succesfully created
// 1 - name is not valid
// 2 - database error
func CreateTag(name string) int {
	//validate name: a-z and '-' only, start&end with a-z.... ^[a-z]([a-z-]*[a-z])?$
	valid, _ := regexp.MatchString("^[a-z]([a-z-]*[a-z])?$", name)
	if !valid {
		return 1
	}
	insertResult, err := TagsCollection.InsertOne(Ctx, bson.D{
		{tagNameField, name},
		{followersField, bson.A{}},
		{eventsField, bson.A{}},
	})
	if err != nil {
		fmt.Println("set: tag.go: CreateTag(: name=", name, err)
		fmt.Println(insertResult)
		return 2
	}
	fmt.Println(insertResult)
	return 0
}

// TagEvent tags the event of eventID with the tag of tagName. email is the email of logged in user
// Returns:
// 0 = successful
// 1 = tag does not exist
// 2 = event does not exist
// 3 = event is not created by the user tagging or not logged in
// 4 = db error on updating TagsCollection
// 5 = db error on updating EventsCollection
func TagEvent(email string, tagName string, eventID primitive.ObjectID) int {
	errorPrefix := "set: tag.go: TagEvent:"
	if _, ok := readTag(tagName, errorPrefix); !ok {
		return 1
	}

	event, ok := readEvent(eventID, errorPrefix)
	if !ok {
		return 2
	}

	if email != event[emailField].(string) {
		return 3
	}

	updateResult, err := TagsCollection.UpdateOne(Ctx, bson.M{tagNameField: tagName}, bson.D{
		{"$push", bson.D{
			{eventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the tag with event", err)
		return 4
	}

	updateResult, err = EventsCollection.UpdateOne(Ctx, bson.M{eventIDField: eventID}, bson.D{
		{"$push", bson.D{
			{tagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the event with tag", err)
		return 5
	}
	return 0
}

// UntagEvent untags the event of eventID with the tag of tagName. email is the email of logged in user
// Returns:
// 0 = successful
// 1 = tag does not exist
// 2 = event does not exist
// 3 = event is not created by the user tagging or not logged in
// 4 = db error on updating TagsCollection
// 5 = db error on updating EventsCollection
func UntagEvent(email string, tagName string, eventID primitive.ObjectID) int {
	errorPrefix := "set: tag.go: UntagEvent:"
	if _, ok := readTag(tagName, errorPrefix); !ok {
		return 1
	}

	event, ok := readEvent(eventID, errorPrefix)
	if !ok {
		return 2
	}

	if email != event[emailField].(string) {
		return 3
	}

	updateResult, err := TagsCollection.UpdateOne(Ctx, bson.M{tagNameField: tagName}, bson.D{
		{"$pull", bson.D{
			{eventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the tag with event", err)
		return 4
	}

	updateResult, err = EventsCollection.UpdateOne(Ctx, bson.M{"_id": eventID}, bson.D{
		{"$pull", bson.D{
			{tagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the event with tag", err)
		return 5
	}
	return 0
}

// FollowTag adds the tag to list of followed tag and adds the user to list of followers
// Returns:
// 0 - success
// 1 - user not found
// 2 - tag not found
// 3 - db update error in usersCollection
// 4 - db update error in tagsCollection
func FollowTag(email string, tagName string) int {
	errorPrefix := "set: tag.go: FollowTag:"
	if _, ok := readUser(email, errorPrefix); !ok {
		return 1
	}
	if _, ok := readTag(tagName, errorPrefix); !ok {
		return 2
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{emailField: email}, bson.D{
		{"$push", bson.D{
			{tagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating userCln", err)
		return 3
	}
	updateResult, err = TagsCollection.UpdateOne(Ctx, bson.M{tagNameField: tagName}, bson.D{
		{"$push", bson.D{
			{followersField, email},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating tagCln", err)
		return 4
	}
	return 0
}

func UnfollowTag(email string, tagName string) int {
	errorPrefix := "set: tag.go: UnfollowTag"
	if _, ok := readUser(email, errorPrefix); !ok {
		return 1
	}
	if _, ok := readTag(tagName, errorPrefix); !ok {
		return 2
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{emailField: email}, bson.D{
		{"$pull", bson.D{
			{tagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating userCln", err)
		return 3
	}
	updateResult, err = TagsCollection.UpdateOne(Ctx, bson.M{tagNameField: tagName}, bson.D{
		{"$pull", bson.D{
			{followersField, email},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating tagCln", err)
		return 4
	}
	return 0
}
func readTag(tagName string, errorPrefix string) (tag bson.M, ok bool) {
	tag = bson.M{}
	err := TagsCollection.FindOne(Ctx, bson.M{
		tagNameField: tagName,
	}).Decode(&tag)
	if err != nil {
		fmt.Println(errorPrefix, "tagName =", tagName, err)
		return nil, false
	}
	return tag, true
}

func readEvent(eventID primitive.ObjectID, errorPrefix string) (event bson.M, ok bool) {
	event = bson.M{}
	err := EventsCollection.FindOne(Ctx, bson.M{
		eventIDField: eventID,
	}).Decode(&event)
	if err != nil {
		fmt.Println(errorPrefix, "eventID =", eventID, err)
		return nil, false
	}
	return event, true
}

func readUser(email string, errorPrefix string) (user bson.M, ok bool) {
	user = bson.M{}
	err := UsersCollection.FindOne(Ctx, bson.M{
		emailField: email,
	}).Decode(&user)
	if err != nil {
		fmt.Println(errorPrefix, "email =", email, err)
		return nil, false
	}
	return user, true
}
