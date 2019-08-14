package ftp

const name = "ftp"

// Config holds the config the ftp option
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
