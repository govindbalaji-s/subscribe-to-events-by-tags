package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"./api"
	"./authzero"
	"./db"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	err := authzero.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	cancelDbCtx := db.Init()
	defer cancelDbCtx()

	rout := mux.NewRouter()

	rout.HandleFunc("/authcallback", authzero.CallbackHandler)
	rout.HandleFunc("/login", authzero.LoginHandler)
	rout.HandleFunc("/logout", authzero.LogoutHandler)
	rout.HandleFunc("/dashboard", authZeroTestHandler)

	tagPostRout := rout.PathPrefix("/api/tag").Methods("POST").Subrouter()
	tagPostRout.HandleFunc("/create/{tagName}", api.CreateTag)
	tagPostRout.HandleFunc("/follow/{tagName}", api.FollowTag)
	tagPostRout.HandleFunc("/unfollow/{tagName}", api.UnfollowTag)
	tagPostRout.HandleFunc("/tag/{tagName}/{eventID}", api.TagEvent)
	tagPostRout.HandleFunc("/untag/{tagName}/{eventID}", api.UntagEvent)

	tagGetRout := rout.PathPrefix("/api/tag").Methods("GET").Subrouter()
	tagGetRout.HandleFunc("/get/{tagName}", api.GetTag)
	tagGetRout.HandleFunc("/search/{query}", api.SearchTags)
	//fmt.Println("going to start")
	log.Fatal(http.ListenAndServe(":8080", rout))
}

func authZeroTestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := authzero.GetAuthSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Check if this is new user and if so , add in database
	CheckUserInDB(session, r)

	temp, _ := template.ParseFiles("testdash.html")
	temp.Execute(w, session.Values["profile"])
	//fmt.Fprintln(w, session.Values["profile"].(map[string]interface{})["nickname"])
}

// CheckUserInDB checks if the logged in user from the session is in db, else he is added.
func CheckUserInDB(session *sessions.Session, r *http.Request) {
	userProfile, ok := session.Values["profile"].(map[string]interface{})
	if !ok {
		//not signed in
		return
	}
	dbResult := &bson.D{{}}
	err := db.UsersCollection.FindOne(db.Ctx, bson.M{
		"email": userProfile["email"]}).Decode(dbResult)
	if err != nil {
		//fmt.Println(err)
		//:::NOTE::::Assuming this error occurs only when no results match
		newUserDoc := bson.D{
			{"email", userProfile["email"]},
			/*{"username", userProfile["username"]},*/
			{"name", ""},
			{"tags", bson.A{}},
			{"archivedEvents", bson.A{}},
			{"subscribedEvents", bson.A{}},
			{"createdEvents", bson.A{}},
			{"queuedPush", bson.A{}},
		}
		insertResult, err := db.UsersCollection.InsertOne(db.Ctx, newUserDoc)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(insertResult)
	}
}
