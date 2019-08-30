package drive

import "fmt"

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
func (p *Provider) Put(file string) {
	fmt.Println("TODO " + file)
}

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {
	return &Provider{}
}
