package models

import "gmail_backup/pkg/storage"

// Settings is the model for settings of this app
type Settings struct {
	StorageOptions []StorageOptions `json:"storage_options"`
}

// StorageOptions ...
type StorageOptions struct {
	Provider storage.Provider `json:"-"`
	Name     string           `json:"name"`
	Active   bool             `json:"active"`
	Path     string           `json:"path"`
	Config   interface{}      `json:"config"`
}
