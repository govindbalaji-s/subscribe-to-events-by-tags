package authzero

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"os"
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
	}
	rawFileContents, _ := ioutil.ReadAll(authConfigFile)
	json.Unmarshal([]byte(rawFileContents), &auth0config)
	authConfigFile.Close()
	Store = sessions.NewFilesystemStore("", []byte(auth0config["sessionKey"].(string)))
	gob.Register(map[string]interface{}{})
	return nil
}
