package storage

// Storage handles the storage options
type Storage struct {
	Providers map[string]Provider
}

// New returns a new Storage
func New() *Storage {
	return &Storage{
		Providers: make(map[string]Provider),
	}
}

// UseProviders adds a list of available providers for use with Storage.
func (s *Storage) UseProviders(viders ...Provider) {
	for _, provider := range viders {
		s.Providers[provider.Name()] = provider
	}
}

// Provider is an interface for the StorageProviders
type Provider interface {
	Name() string
}

// Register a new Provider in storage map
func (s *Storage) Register(p Provider) {
	s.Providers[p.Name()] = p
}

// GetProviders returns the registered providers
func (s *Storage) GetProviders() map[string]Provider {
	return s.Providers
}
