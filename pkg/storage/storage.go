package storage

// Storage holds the different storage options
type Storage struct{}

// Types returns the storage types
var Types = map[string]string{
	"ftp":          "ftp",
	"google_drive": "Google Drive",
}
