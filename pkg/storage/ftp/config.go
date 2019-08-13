package ftp

const name = "ftp"

// Provider implements storage.Provider for the ftp file storage.
type Provider struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
}

// Name returns ftp
func (p *Provider) Name() string {
	return name
}

// New initializer for Provider struct ftp
func New(username, password, path string) *Provider {
	return &Provider{
		Username: username,
		Password: password,
		Path:     path,
	}
}
