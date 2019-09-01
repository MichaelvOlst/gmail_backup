package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HandlerBackupAccount backups an account and returns a websocket connection
func (a *API) HandlerBackupAccount(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	aid, ok := vars["id"]
	if !ok {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: "Id not valid"})
	}
	id, _ := strconv.Atoi(aid)

	ac, err := a.db.GetAccountByID(id)
	if err != nil {
		return respond(w, http.StatusNotFound, envelope{Error: "Could find an account with this Id"})
	}

	err = a.gmail.Backup(ac)
	if err != nil {
		return respond(w, http.StatusUnprocessableEntity, envelope{Error: fmt.Sprintf("Unable to authenticate: %v", err)})
	}

	return respond(w, http.StatusOK, envelope{Result: ac})
}
