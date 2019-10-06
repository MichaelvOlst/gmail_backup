package models

import (
	"gmail_backup/pkg/storage/drive"
	"gmail_backup/pkg/storage/dropbox"
	"gmail_backup/pkg/storage/ftp"
)

// Settings is the model for settings of this app
type Settings struct {
	StorageOptions StorageOptions `json:"storage_options"`
}

// StorageOptions all the storage options the app provides
type StorageOptions struct {
	Dropbox     Dropbox     `json:"dropbox"`
	Ftp         Ftp         `json:"ftp"`
	GoogleDrive GoogleDrive `json:"google_drive"`
}

// StorageOption ...
type StorageOption struct {
	Option string `json:"option"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// Dropbox config for the dropbox storage option
type Dropbox struct {
	StorageOption StorageOption
	Config        dropbox.Config
}

// Ftp config for the ftp storage option
type Ftp struct {
	StorageOption StorageOption
	Config        ftp.Config
}

// GoogleDrive config for the ftp storage option
type GoogleDrive struct {
	StorageOption StorageOption
	Config        drive.Config
}
