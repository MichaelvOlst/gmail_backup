package database

import (
	"gmail_backup/pkg/config"

	"github.com/asdine/storm"
)

// Store returns a db object
type Store struct {
	*storm.DB
}

// New returns a storm.db object
func New(cfg *config.Config) (*Store, error) {
	db, err := storm.Open(cfg.Database.Filename)
	if err != nil {
		return nil, err
	}
	return &Store{db}, nil
}
