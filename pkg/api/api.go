package api

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/sessions"
)

// API ...
type API struct {
	config   *config.Config
	db       *database.Store
	box      *packr.Box
	sessions sessions.Store
}

// New returns a new API
func New(config *config.Config, db *database.Store, box *packr.Box) *API {
	return &API{
		config:   config,
		db:       db,
		box:      box,
		sessions: sessions.NewCookieStore([]byte(config.Server.Secret)),
	}
}
