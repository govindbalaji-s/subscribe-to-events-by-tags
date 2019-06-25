package authzero

import (
	"net/http"
	"net/url"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	domain := auth0config["domain"].(string)

	var Url *url.URL
	Url, err := url.Parse("https://" + domain)

	if err != nil {
		panic("boom")
	}

	Url.Path += "/v2/logout"
	parameters := url.Values{}
	parameters.Add("returnTo", postLogoutURL)
	parameters.Add("client_id", auth0config["ClientID"].(string))
	Url.RawQuery = parameters.Encode()

	http.Redirect(w, r, Url.String(), http.StatusTemporaryRedirect)
}
