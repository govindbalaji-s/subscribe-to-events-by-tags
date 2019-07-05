package authzero

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var (
	Store       *sessions.FilesystemStore
	auth0config map[string]interface{}
)

const (
	authConfigFilePath = "./authzero/auth0config.json"
	authCallbackURL    = "http://127.0.0.1:8080/authcallback"
	postAuthPath       = "/dashboard"
	postLogoutURL      = "http://127.0.0.1:8080/dashboard"
)

func Init() error {
	authConfigFile, err := os.Open(authConfigFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	rawFileContents, _ := ioutil.ReadAll(authConfigFile)
	json.Unmarshal([]byte(rawFileContents), &auth0config)
	authConfigFile.Close()
	Store = sessions.NewFilesystemStore("", []byte(auth0config["sessionKey"].(string)))
	gob.Register(map[string]interface{}{})
	return nil
}

func GetAuthSession(r *http.Request) (session *sessions.Session, err error) {
	session, err = Store.Get(r, "auth-session")
	if err != nil {
		fmt.Println("Errrrrrrrrrrrrrrrrrrrrrrrrr", err)
	}
	return
}

func GetUserEmailFromSession(r *http.Request) (string, error) {
	session, err := GetAuthSession(r)
	if err != nil {
		return "", err
	}

	userProfile, ok := session.Values["profile"].(map[string]interface{})
	if !ok {
		return "", NotSignedInError{}
	}
	return userProfile["email"].(string), nil
}

type NotSignedInError struct{}

func (err NotSignedInError) Error() string {
	return "No user signed in."
}
