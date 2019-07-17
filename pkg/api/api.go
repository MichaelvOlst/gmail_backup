package api

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"

	"github.com/gorilla/sessions"
)

// API ...
type API struct {
	cfg      *config.Config
	sessions sessions.Store
	db       *database.Store
}

// New returns a new API
func New(config *config.Config, db *database.Store) *API {

	return &API{
		cfg:      config,
		db:       db,
		sessions: sessions.NewCookieStore([]byte(config.Server.Secret)),
	}
}
