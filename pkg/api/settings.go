package api

import (
	"fmt"
	"gmail_backup/pkg/models"
	"net/http"
)

// HandlerGetSettings gets the app settings
func (a *API) HandlerGetSettings(w http.ResponseWriter, r *http.Request) error {

	s := models.Settings{}
	err := a.db.Get("settings", "settings", &s)
	if err != nil {
		fmt.Println(err)
		return respond(w, http.StatusInternalServerError, envelope{Error: err})
	}

	return respond(w, http.StatusOK, envelope{Result: s})
}
