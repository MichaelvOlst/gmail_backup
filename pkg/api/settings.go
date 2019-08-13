package api

import (
	"net/http"
)

// HandlerGetSettings gets the app settings
func (a *API) HandlerGetSettings(w http.ResponseWriter, r *http.Request) error {

	// a.storage.getProviders()

	return respond(w, http.StatusOK, envelope{Result: "getting settings"})
}
