package dropbox

const name = "dropbox"

// Config holds the config for the Dropbox option
type Config struct {
	AccessToken string `json:"accesstoken"`
}

// Provider implements storage.Provider for the ftp file storage.
type Provider struct {
}

// Name returns ftp
func (p *Provider) Name() string {
	return name
}

// New initializer for Provider struct ftp
func New() *Provider {
	return &Provider{}
}
