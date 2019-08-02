package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asdine/storm"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *login) Sanitize() {
	l.Email = strings.ToLower(strings.TrimSpace(l.Email))
}

// HandlerLogin starts the server
func (a *API) HandlerLogin(w http.ResponseWriter, r *http.Request) error {
	var l login
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		return err
	}
	l.Sanitize()

	u, err := a.db.GetUserByEmail(l.Email)
	if err != nil && err != storm.ErrNotFound {
		return err
	}

	if err == storm.ErrNotFound || u.ComparePassword(l.Password) != nil {
		return respond(w, http.StatusUnauthorized, envelope{Error: "invalid_credentials"})
	}

	session, _ := a.sessions.Get(r, "auth")
	session.Values["user_id"] = u.ID
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return respond(w, http.StatusOK, envelope{Result: true})
}

// HandlerLogout starts the server
func (a *API) HandlerLogout(w http.ResponseWriter, r *http.Request) error {
	session, _ := a.sessions.Get(r, "auth")
	if !session.IsNew {
		session.Options.MaxAge = -1
		err := session.Save(r, w)
		if err != nil {
			return err
		}
	}
	return respond(w, http.StatusOK, envelope{Result: true})
}

// HandlerSession checks if the user is still logged in
func (a *API) HandlerSession(w http.ResponseWriter, r *http.Request) error {

	session, _ := a.sessions.Get(r, "auth")
	if !session.IsNew {
		return respond(w, http.StatusOK, envelope{Result: true})
	}

	return respond(w, http.StatusOK, envelope{Result: false})
}

// Authorize is middleware that aborts the request if unauthorized
func (a *API) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := a.sessions.Get(r, "auth")
		// an err is returned if cookie has been tampered with, so check that
		if err != nil {
			respond(w, http.StatusUnauthorized, envelope{Error: "unauthorized"})
			return
		}

		userID, ok := session.Values["user_id"]
		if session.IsNew || !ok {
			respond(w, http.StatusUnauthorized, envelope{Error: "unauthorized"})
			return
		}

		// validate user ID in session
		if _, err := a.db.GetUserByID(userID.(int)); err != nil {
			respond(w, http.StatusUnauthorized, envelope{Error: "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
