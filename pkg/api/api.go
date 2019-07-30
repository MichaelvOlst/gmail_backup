package api

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"

	"github.com/gobuffalo/packr"
)

// API ...
type API struct {
	config *config.Config
	db     *database.Store
	box    *packr.Box
}

// New returns a new API
func New(config *config.Config, db *database.Store, box *packr.Box) *API {
	return &API{
		config: config,
		db:     db,
		box:    box,
	}
}
