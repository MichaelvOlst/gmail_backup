package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type envelope struct {
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// Handler is our custom HTTP handler with error returns
type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		HandleError(w, r, err)
	}
}

// HandleError handles errors
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	logrus.WithFields(logrus.Fields{
		"request": r.Method + " " + r.RequestURI,
		"error":   err,
	}).Error("error handling request")

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("false"))
}

func respond(w http.ResponseWriter, statusCode int, d interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(d)
	return err
}

// HandlerFunc takes a custom Handler func and converts it to http.HandlerFunc
func HandlerFunc(fn Handler) http.HandlerFunc {
	return http.HandlerFunc(Handler(fn).ServeHTTP)
}

// ServeFileHandler serves the file to the browser
func (a *API) ServeFileHandler(filename string) http.Handler {
	return HandlerFunc(a.serveFile(filename))
}

func (a *API) serveFile(filename string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		f, err := a.box.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			return err
		}

		// setting security and cache headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Xss-Protection", "1; mode=block")
		w.Header().Set("Cache-Control", "max-age=432000") // 5 days

		http.ServeContent(w, r, filename, d.ModTime(), f)
		return nil
	}
}

// NotFoundHandler will return a 404 page
func (a *API) NotFoundHandler() http.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusNotFound)
		w.Write(a.box.Bytes("404.html"))
		return nil
	})
}
