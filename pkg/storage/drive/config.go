package drive

const name = "google_drive"

// Config holds the config the ftp option
type Config struct {
	Client string `json:"client"`
	Secret string `json:"secret"`
}

// Provider implements storage.Provider for the ftp file storage.
type Provider struct {
	Config Config
}

// Name returns ftp
func (p *Provider) Name() string {
	return name
}

// New initializer for Provider struct ftp
func New() *Provider {
	return &Provider{}
}
