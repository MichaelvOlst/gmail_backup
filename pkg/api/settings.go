package api

import (
	"encoding/json"
	"gmail_backup/pkg/models"
	"net/http"
)

// HandlerGetSettings gets the app settings
func (a *API) HandlerGetSettings(w http.ResponseWriter, r *http.Request) error {

	s := models.Settings{}
	err := a.db.Get("settings", "settings", &s)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	return respond(w, http.StatusOK, envelope{Result: s})
}

// HandlerSaveSettings gets the app settings
func (a *API) HandlerSaveSettings(w http.ResponseWriter, r *http.Request) error {
	var s models.Settings

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		return err
	}

	err = a.db.Set("settings", "settings", &s)
	if err != nil {
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	a.storage.ClearProviders()
	// for _, val := range s.StorageOptions {
	// 	if val.Active {
	// 		a.storage.Register(val.Option, val.Config)
	// 	}
	// }

	return respond(w, http.StatusOK, envelope{Result: s})
}
