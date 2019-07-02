package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"set/authzero"
	"set/db"
	"set/util"
	"strconv"

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret := db.CreateTag(vars["tagName"])
	switch ret {
	case 0:
		writeSuccess(w)
	case 1:
		writeFailed(w, "no user signed in", http.StatusBadRequest)
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
		writeSuccess(w)
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
		writeSuccess(w)
	case ret == 1:
		writeFailed(w, "user is not creator", http.StatusBadRequest)
	case ret == 2:
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
	followers := tag[db.FollowersField].(primitive.A)
	events := tag[db.EventsField].(primitive.A)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(map[string]interface{}{
		"result": "success",
		"error":  "",
		"data": map[string]string{
			"tagName":       tagName,
			"isFollowing":   strconv.FormatBool(util.Contains(followers, userEmail)),
			"noOfFollowers": string(len(followers)),
			"noOfEvents":    string(len(events)),
		},
	})
	w.Write(jsonResp)
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
	var data []map[string]string
	for _, tag := range tags {
		followers := tag[db.FollowersField].(primitive.A)
		events := tag[db.EventsField].(primitive.A)
		data = append(data, map[string]string{
			"tagName":       tag[db.TagNameField].(string),
			"isFollowing":   strconv.FormatBool(util.Contains(followers, userEmail)),
			"noOfFollowers": string(len(followers)),
			"noOfEvents":    string(len(events)),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(map[string]interface{}{
		"result": "success",
		"error":  "",
		"data":   data,
	})
	w.Write(jsonResp)
}
func writeSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(map[string]string{
		"result": "success",
		"error":  "",
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
