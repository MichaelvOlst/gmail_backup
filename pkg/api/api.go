package api

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/storage"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/sessions"
)

// API ...
type API struct {
	config   *config.Config
	db       *database.Store
	box      *packr.Box
	sessions sessions.Store
	storage  *storage.Storage
}

// New returns a new API
func New(config *config.Config, db *database.Store, box *packr.Box, s *storage.Storage) *API {
	return &API{
		config:   config,
		db:       db,
		box:      box,
		sessions: sessions.NewCookieStore([]byte(config.Server.Secret)),
		storage:  s,
	}
}
