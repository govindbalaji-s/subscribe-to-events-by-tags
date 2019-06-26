package main

import (
	"./authzero"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	err := authzero.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

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
	temp, _ := template.ParseFiles("testdash.html")
	temp.Execute(w, session.Values["profile"])
}
