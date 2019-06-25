package main

import (
	"./authzero"
	"fmt"
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
	//fmt.Println("going to start")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
