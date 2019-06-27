package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"./authzero"
	"./db"
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

	http.HandleFunc("/authcallback", authzero.CallbackHandler)
	http.HandleFunc("/login", authzero.LoginHandler)
	http.HandleFunc("/logout", authzero.LogoutHandler)
	http.HandleFunc("/dashboard", authZeroTestHandler)
	//fmt.Println("going to start")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authZeroTestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := authzero.Store.Get(r, "auth-session")
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

func CheckUserInDB(session *sessions.Session, r *http.Request) {
	userProfile := session.Values["profile"].(map[string]interface{})
	usersCollection := db.Db.Collection("users")
	dbResult := &bson.D{{}}
	err := usersCollection.FindOne(db.Ctx, bson.M{
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
		insertResult, err := usersCollection.InsertOne(db.Ctx, newUserDoc)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(insertResult)
	}
}
