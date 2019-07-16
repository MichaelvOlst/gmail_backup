package api

import (
	"gmail_backup/pkg/config"

	"github.com/gorilla/sessions"
)

// API ...
type API struct {
	cfg      *config.Config
	sessions sessions.Store
}

// New returns a new API
func New(config *config.Config) *API {

	return &API{
		cfg:      config,
		sessions: sessions.NewCookieStore([]byte(config.Server.Secret)),
	}
}
