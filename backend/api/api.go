package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"set/backend/authzero"
	"set/backend/db"
	"set/backend/util"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	err := authzero.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	cancelDbCtx := db.Init()
	defer cancelDbCtx()
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		fmt.Println("Errrrrrrrrr", err)
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	tagName := vars["tagName"]
	//validate name: a-z and '-' only, start&end with a-z.... ^[a-z]([a-z-]*[a-z])?$
	if !IsValidTagName(tagName) {
		writeFailed(w, "invalid tag name", http.StatusBadRequest)
		return
	}
	ret := db.CreateTag(tagName)
	switch ret {
	case 0:
		writeSuccess(w, nil)
	case 2:
		writeFailed(w, "tag already exists", http.StatusBadRequest)
	case 3:
		writeFailed(w, "database error", http.StatusInternalServerError)
	}
}

func followTagUtil(w http.ResponseWriter, r *http.Request, toFollow bool) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	ret := db.FollowTag(userEmail, vars["tagName"], toFollow)
	switch true {
	case ret == 0:
		writeSuccess(w, nil)
	case ret == 1:
		writeFailed(w, "user dne", http.StatusBadRequest)
	case ret == 2:
		writeFailed(w, "tag dne", http.StatusBadRequest)
	case ret == 3 || ret == 4:
		writeFailed(w, "database erro", http.StatusInternalServerError)
	}
}
func FollowTag(w http.ResponseWriter, r *http.Request) {
	followTagUtil(w, r, true)
}

func UnfollowTag(w http.ResponseWriter, r *http.Request) {
	followTagUtil(w, r, false)
}

func tagEventUtil(w http.ResponseWriter, r *http.Request, toTag bool) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	eventID, err := primitive.ObjectIDFromHex(vars["eventID"])
	if err != nil {
		writeFailed(w, "not an objectid", http.StatusBadRequest)
		return
	}
	ret := db.TagEvent(userEmail, vars["tagName"], eventID, toTag)
	switch true {
	case ret == 0:
		writeSuccess(w, nil)
	case ret == 1:
		writeFailed(w, "tag does not exist", http.StatusBadRequest)
	case ret == 2:
		writeFailed(w, "event dne", http.StatusBadRequest)
	case ret == 3:
		writeFailed(w, "user is not cretor", http.StatusBadRequest)
	case ret == 4 || ret == 5:
		writeFailed(w, "database error", http.StatusInternalServerError)
	}
}

func TagEvent(w http.ResponseWriter, r *http.Request) {
	tagEventUtil(w, r, true)
}
func UntagEvent(w http.ResponseWriter, r *http.Request) {
	tagEventUtil(w, r, false)
}

func GetTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	tagName := vars["tagName"]
	tag, ok := db.ReadTag(tagName, "set: api/tag.go: GetTag")
	if !ok {
		writeFailed(w, "tag dne", http.StatusBadRequest)
		return
	}
	data := tagToSendable(tag, userEmail)
	writeSuccess(w, data)
}

func SearchTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	query := vars["query"]
	tags, ok := db.SearchTags(query)
	if !ok {
		writeFailed(w, "database error", http.StatusInternalServerError)
		return
	}
	var data []map[string]interface{}
	for _, tag := range tags {
		data = append(data, tagToSendable(tag, userEmail))
	}
	writeSuccess(w, data)
}

func tagToSendable(tag bson.M, userEmail string) map[string]interface{} {
	followers := tag[db.FollowersField].(primitive.A)
	events := tag[db.EventsField].(primitive.A)
	tagName := tag[db.TagNameField].(string)
	data := map[string]interface{}{
		"tagName":       tagName,
		"isFollowing":   strconv.FormatBool(util.Contains(followers, userEmail)),
		"noOfFollowers": strconv.Itoa(len(followers)),
		"noOfEvents":    strconv.Itoa(len(events)),
	}
	return data
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeFailed(w, "failed to read Body", http.StatusInternalServerError)
		return
	}
	var reqVars map[string]interface{}
	err = json.Unmarshal(body, &reqVars)
	if err != nil {
		writeFailed(w, "failed to parse Body", http.StatusInternalServerError)
		return
	}
	//validate request body
	//name, venue, time, duration, tags
	reqvar, allOk := reqVars["name"]
	eventName, ok := reqvar.(string)
	allOk = allOk && ok

	reqvar, ok = reqVars["venue"]
	allOk = allOk && ok
	eventVenue, ok := reqvar.(string)
	allOk = allOk && ok

	reqvar, ok = reqVars["time"]
	allOk = allOk && ok
	eventTimeInString, ok := reqvar.(string)
	allOk = allOk && ok
	eventTime, err := strconv.ParseInt(eventTimeInString, 10, 64)
	allOk = allOk && (err == nil)

	reqvar, ok = reqVars["duration"]
	allOk = allOk && ok
	eventDurationAsString, ok := reqvar.(string)
	allOk = allOk && ok
	eventDuration, err := strconv.ParseInt(eventDurationAsString, 10, 64)
	allOk = allOk && (err == nil)

	reqvar, ok = reqVars["tags"]
	allOk = allOk && ok
	eventTagsGen, ok := reqvar.([]interface{})
	var eventTags []string
	for _, t := range eventTagsGen {
		tagName, ok := t.(string)
		allOk = allOk && ok
		eventTags = append(eventTags, tagName)
	}
	allOk = allOk && ok
	for _, tagName := range eventTags {
		allOk = allOk && IsValidTagName(tagName)
	}
	if !allOk {
		writeFailed(w, "invalid args", http.StatusBadRequest)
		return
	}

	ret, eventID := db.CreateEvent(eventName, eventVenue, eventTime, eventDuration, eventTags, userEmail)

	switch true {
	case ret == 0:
		writeSuccess(w, map[string]interface{}{
			"eventID": eventID.Hex(),
		})
	case ret == 1:
		writeFailed(w, "signed in user not found.", http.StatusBadRequest)
	case ret == 2 || ret == 3 || ret == 4:
		fmt.Println("CreateEvent failed ret=", ret)
		writeFailed(w, "database error", http.StatusInternalServerError)
	}

}

func EditEvent(w http.ResponseWriter, r *http.Request) {
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeFailed(w, "failed to read Body", http.StatusInternalServerError)
		return
	}
	var reqVars map[string]interface{}
	err = json.Unmarshal(body, &reqVars)
	if err != nil {
		writeFailed(w, "failed to parse Body", http.StatusInternalServerError)
		return
	}

	eventDetailsMap := bson.M{}

	reqvar, ok := reqVars["eventID"]
	if !ok {
		writeFailed(w, "event id needed.", http.StatusBadRequest)
		return
	}
	eventIDAsString, ok := reqvar.(string)
	if !ok {
		writeFailed(w, "event id can not be parsed", http.StatusBadRequest)
		return
	}
	eventID, err := primitive.ObjectIDFromHex(eventIDAsString)
	if err != nil {
		writeFailed(w, "event id is not valid hex", http.StatusBadRequest)
		return
	}

	reqvar, ok = reqVars["name"]
	if ok {
		eventName, ok := reqvar.(string)
		if ok {
			eventDetailsMap[db.EventNameField] = eventName
		}
	}

	reqvar, ok = reqVars["venue"]
	if ok {
		eventVenue, ok := reqvar.(string)
		if ok {
			eventDetailsMap[db.EventVenueField] = eventVenue
		}
	}

	reqvar, ok = reqVars["time"]
	if ok {
		eventTimeInString, ok := reqvar.(string)
		if ok {
			eventTime, err := strconv.ParseInt(eventTimeInString, 10, 64)
			if err == nil {
				eventDetailsMap[db.EventTimeField] = eventTime
			}
		}
	}

	reqvar, ok = reqVars["duration"]
	if ok {
		eventDurationInString, ok := reqvar.(string)
		if ok {
			eventDuration, err := strconv.ParseInt(eventDurationInString, 10, 64)
			if err == nil {
				eventDetailsMap[db.EventDurationField] = eventDuration
			}
		}
	}

	ret := db.EditEventDetails(eventID, eventDetailsMap, userEmail)

	switch true {
	case ret == 0:
		writeSuccess(w, nil)
	case ret == 1:
		writeFailed(w, "event dne", http.StatusBadRequest)
	case ret == 2:
		writeFailed(w, "only creator can edit", http.StatusBadRequest)
	case ret == 3:
		writeFailed(w, "database error", http.StatusInternalServerError)
	}
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	eventID, err := primitive.ObjectIDFromHex(vars["eventID"])
	if err != nil {
		writeFailed(w, "event id not valid", http.StatusBadRequest)
		return
	}
	ret := db.DeleteEvent(eventID, userEmail)

	switch true {
	case ret == 0:
		writeSuccess(w, nil)
	case ret == 1:
		writeFailed(w, "event dne", http.StatusBadRequest)
	case ret == 2:
		writeFailed(w, "only creator can delete", http.StatusBadRequest)
	case ret == 3 || ret == 4 || ret == 5 || ret == 6:
		writeFailed(w, "database error", http.StatusInternalServerError)
	}
}

func subscribeToEventUtil(w http.ResponseWriter, r *http.Request, toSubscribe bool) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	eventID, err := primitive.ObjectIDFromHex(vars["eventID"])
	if err != nil {
		writeFailed(w, "event id not valid", http.StatusBadRequest)
		return
	}

	ret := db.SubscribeToEvent(eventID, userEmail, toSubscribe)

	switch true {
	case ret == 0:
		writeSuccess(w, nil)
	case ret == 1:
		writeFailed(w, "event dne", http.StatusBadRequest)
	case ret == 2:
		writeFailed(w, "user dne", http.StatusBadRequest)
	case ret == 3 || ret == 4:
		writeFailed(w, "database error", http.StatusInternalServerError)
	}
}
func SubscribeToEvent(w http.ResponseWriter, r *http.Request) {
	subscribeToEventUtil(w, r, true)
}

func UnsubscribeToEvent(w http.ResponseWriter, r *http.Request) {
	subscribeToEventUtil(w, r, false)
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	eventID, err := primitive.ObjectIDFromHex(vars["eventID"])
	if err != nil {
		writeFailed(w, "event id not valid", http.StatusBadRequest)
		return
	}
	//name, venue, time, duration, tags, (creator name?), subscriber count, isSubscribed
	event, ok := db.ReadEvent(eventID, "set: api.go GetEvent")
	if !ok {
		writeFailed(w, "event dne", http.StatusBadRequest)
		return
	}
	data := eventToSendable(event, userEmail)
	writeSuccess(w, data)
}

func SearchEventsByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail, err := authzero.GetUserEmailFromSession(r)
	if err != nil {
		writeFailed(w, "no user signed in", http.StatusBadRequest)
		return
	}
	query := vars["query"]
	events, ok := db.SearchEventsByName(query)
	if !ok {
		writeFailed(w, "database error", http.StatusInternalServerError)
		return
	}
	var data []map[string]interface{}
	for _, event := range events {
		data = append(data, eventToSendable(event, userEmail))
	}
	writeSuccess(w, data)
}

func eventToSendable(event bson.M, userEmail string) map[string]interface{} {
	return map[string]interface{}{
		"eventID":        event[db.EventIDField].(primitive.ObjectID).Hex(),
		"name":           event[db.EventNameField].(string),
		"venue":          event[db.EventVenueField].(string),
		"time":           strconv.FormatInt(event[db.EventTimeField].(int64), 10),
		"duration":       strconv.FormatInt(event[db.EventDurationField].(int64), 10),
		"tags":           event[db.EventTagsField].(primitive.A),
		"noOfSubsribers": strconv.Itoa(len(event[db.EventSubscribersField].(primitive.A))),
		"isSubscribed":   util.Contains(event[db.EventSubscribersField].(primitive.A), userEmail),
	}
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(map[string]interface{}{
		"result": "success",
		"error":  "",
		"data":   data,
	})
	w.Write(jsonResp)
}

func writeFailed(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonResp, _ := json.Marshal(map[string]string{
		"result": "failed",
		"error":  err,
	})
	w.Write(jsonResp)
}

func IsValidTagName(tagName string) (valid bool) {
	valid, _ = regexp.MatchString("^[a-z]([a-z-]*[a-z])?$", tagName)
	return valid
}
