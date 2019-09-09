package api

import (
	"fmt"
	"gmail_backup/pkg/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandlerBackupAccount backups an account and returns a websocket connection
func (a *API) HandlerBackupAccount(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}

	vars := mux.Vars(r)
	aid, ok := vars["id"]
	if !ok {
		ws.WriteJSON(envelope{Error: "Could not parse id"})
		return
	}
	id, _ := strconv.Atoi(aid)

	ac, err := a.db.GetAccountByID(id)
	if err != nil {
		ws.WriteJSON(envelope{Error: "Could not find an account with this Id"})
		return
	}

	go a.gmail.Backup(ac, a.storage)

	go a.writer(ws, ac)

	// ac, err := a.db.GetAccountByID(id)
	// if err != nil {
	// 	return respond(w, http.StatusNotFound, envelope{Error: "Could find an account with this Id"})
	// }

	// ok, err = a.gmail.Backup(ac)
	// if err != nil {
	// 	return respond(w, http.StatusUnprocessableEntity, envelope{Error: fmt.Sprintf("Unable to authenticate: %v", err)})
	// }

	// _ = ok

	// return respond(w, http.StatusOK, envelope{Result: ac})
}

func (a *API) writer(conn *websocket.Conn, ac *models.Account) {

	for {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {

			ab, err := a.db.GetAccountByID(ac.ID)
			if err != nil {
				fmt.Println(err)
				return
			}

			// and finally we write this JSON string to our WebSocket
			// connection and record any errors if there has been any
			if err := conn.WriteJSON(ab); err != nil {
				fmt.Println(err)
				conn.Close()
				return
			}
		}
	}
}
