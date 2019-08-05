package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes starts the server
func (a *API) Routes() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/auth/login", HandlerFunc(a.HandlerLogin)).Methods(http.MethodPost)
	r.Handle("/auth/logout", HandlerFunc(a.HandlerLogout)).Methods(http.MethodPost)
	r.Handle("/auth/session", HandlerFunc(a.HandlerSession)).Methods(http.MethodGet)

	r.Handle("/api/google-url", a.Authorize(HandlerFunc(a.HandlerGetGoogleURL))).Methods(http.MethodGet)

	r.Handle("/api/accounts", a.Authorize(HandlerFunc(a.HandlerGetAllAccounts))).Methods(http.MethodGet)
	r.Handle("/api/accounts/{id:[0-9]+}", a.Authorize(HandlerFunc(a.HandlerGetSingleAccount))).Methods(http.MethodGet)
	r.Handle("/api/accounts", a.Authorize(HandlerFunc(a.HandlerCreateAccount))).Methods(http.MethodPost)
	r.Handle("/api/accounts/{id:[0-9]+}", a.Authorize(HandlerFunc(a.HandlerUpdateAccount))).Methods(http.MethodPatch)
	r.Handle("/api/accounts/{id:[0-9]+}", a.Authorize(HandlerFunc(a.HandlerDeleteAccount))).Methods(http.MethodDelete)

	spa := spaHandler{staticPath: "public/", indexPath: "index.html", box: a.box}
	r.PathPrefix("/").Handler(spa)

	r.Path("/").Handler(a.ServeFileHandler("index.html"))
	// r.Path("/index.html").Handler(a.ServeFileHandler("index.html"))
	// r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(a.box)))
	r.NotFoundHandler = a.NotFoundHandler()

	return r
}
