package drive

import (
	"fmt"
	"os"
)

const name = "google_drive"

// Config holds the config the ftp option
type Config struct {
	Client string `json:"client"`
	Secret string `json:"secret"`
}

// Provider implements storage.Provider for the ftp file storage.
type Provider struct {
}

// Name returns google_drive
func (p *Provider) Name() string {
	return name
}

// ListFolder returns google_drive
func (p *Provider) ListFolder() {
	fmt.Println("TODO")
}

// Put returns google_drive
func (p *Provider) Put(filename, path string, file *os.File) error {
	fmt.Println("TODO " + filename)
	return nil
}

// Mkdir returns google_drive
func (p *Provider) Mkdir(path string) error {
	fmt.Println("TODO " + path)
	return nil
}

// IsNotExists check if a folder already exists
func (p *Provider) IsNotExists(err error) bool {
	return false
	// cerr, ok := err.(files.CreateFolderAPIError)
	// if !ok {
	// 	return false
	// }

	// if cerr.APIError.Error() == "path/conflict/folder/" {
	// 	return true
	// }
}

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {
	return &Provider{}
}
