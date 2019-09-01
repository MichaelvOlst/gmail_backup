package api

import (
	"net/http"
)

// HandlerGetGoogleURL gets the URL for the token
func (a *API) HandlerGetGoogleURL(w http.ResponseWriter, r *http.Request) error {

	authURL := a.gmail.GetAuthCodeURL()

	return respond(w, http.StatusOK, envelope{Result: authURL})
}
