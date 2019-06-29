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
	if _, ok := readTag(tagName, errorPrefix); !ok {
		return 1
	}

	event, ok := readEvent(eventID, errorPrefix)
	if !ok {
		return 2
	}

	if email != event[eventCreatorField].(string) {
		return 3
	}

	updateResult, err := TagsCollection.UpdateOne(Ctx, bson.M{tagNameField: tagName}, bson.D{
		{updateOp, bson.D{
			{eventsField, eventID},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating the tag with event", err)
		return 4
	}

	updateResult, err = EventsCollection.UpdateOne(Ctx, bson.M{eventIDField: eventID}, bson.D{
		{updateOp, bson.D{
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
	if _, ok := readUser(email, errorPrefix); !ok {
		return 1
	}
	if _, ok := readTag(tagName, errorPrefix); !ok {
		return 2
	}
	updateResult, err := UsersCollection.UpdateOne(Ctx, bson.M{emailField: email}, bson.D{
		{updateOp, bson.D{
			{tagsField, tagName},
		}},
	})
	fmt.Println(updateResult)
	if err != nil {
		fmt.Println(errorPrefix+" on updating userCln", err)
		return 3
	}
	updateResult, err = TagsCollection.UpdateOne(Ctx, bson.M{tagNameField: tagName}, bson.D{
		{updateOp, bson.D{
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
