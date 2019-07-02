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
// 2 - tag already exists
// 3 - database error
func CreateTag(name string) int {
	//validate name: a-z and '-' only, start&end with a-z.... ^[a-z]([a-z-]*[a-z])?$
	errorPrefix := "set: db/tag.go: CreateTag:"
	valid, _ := regexp.MatchString("^[a-z]([a-z-]*[a-z])?$", name)
	if !valid {
		return 1
	}
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
