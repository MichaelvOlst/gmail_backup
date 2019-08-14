package api

import (
	"gmail_backup/pkg/models"
	"net/http"
)

// HandlerGetSettings gets the app settings
func (a *API) HandlerGetSettings(w http.ResponseWriter, r *http.Request) error {

	// a.storage.getProviders()
	var s models.Settings
	err := a.db.Get("settings", "settings", &s)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Result: "Could not retrieve settings"})
	}

	return respond(w, http.StatusOK, envelope{Result: s})
}
