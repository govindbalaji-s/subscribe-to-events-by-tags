package authzero

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/oauth2"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	domain := auth0config["domain"].(string)
	aud := ""

	conf := &oauth2.Config{
		ClientID:     auth0config["ClientID"].(string),
		ClientSecret: auth0config["ClientSecret"].(string),
		RedirectURL:  authCallbackURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	if aud == "" {
		aud = "https://" + domain + "/userinfo"
	}

	// Generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := Store.Get(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", aud)
	url := conf.AuthCodeURL(state, audience)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
