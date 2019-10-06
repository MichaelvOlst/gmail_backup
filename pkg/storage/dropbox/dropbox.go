package dropbox

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

const name = "dropbox"
const chunkSize int64 = 1 << 24

// Config holds the config for the Dropbox option
type Config struct {
	AccessToken string `json:"accesstoken"`
}

// Provider implements storage.Provider for the dropbox file storage.
type Provider struct {
	client files.Client
}

// Name returns dropbox
func (p *Provider) Name() string {
	return name
}

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {
	config := dropbox.Config{
		Token: cfg.AccessToken,
	}

	c := files.New(config)

	return &Provider{c}
}
