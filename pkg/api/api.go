package api

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
)

// API ...
type API struct {
	config *config.Config
	db     *database.Store
}

// New returns a new API
func New(config *config.Config, db *database.Store) *API {
	return &API{
		config: config,
		db:     db,
	}
}
