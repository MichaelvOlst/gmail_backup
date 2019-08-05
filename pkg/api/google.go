package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// HandlerGetGoogleURL gets the URL for the accesstoken
func (a *API) HandlerGetGoogleURL(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadFile(a.config.Google.File)
	if err != nil {
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	token := getTokenFromWeb(config)
	return respond(w, http.StatusOK, envelope{Result: token})
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}
