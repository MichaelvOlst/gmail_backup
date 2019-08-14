package models

import "gmail_backup/pkg/storage"

// Settings is the model for settings of this app
type Settings struct {
	StorageOptions []StorageOptions `json:"storage_options"`
}

// StorageOptions ...
type StorageOptions struct {
	Option storage.Provider `json:"option"`
	Active bool             `json:"active"`
	Path   string           `json:"path"`
	Config interface{}      `json:"config"`
}
