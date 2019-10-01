package storage

import (
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage/drive"
	"gmail_backup/pkg/storage/dropbox"
	"gmail_backup/pkg/storage/ftp"
	"io"
	"os"

	"github.com/pkg/errors"
)

// Storage handles the storage options
type Storage struct {
	Providers map[string]Provider
}

// Config default config
type Config interface {
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
	ListFolder()
	Put(filename, path string, file *os.File, r io.Reader) error
	Mkdir(path string) error
	IsNotExists(err error) bool
}

// Register a new Provider in storage map
func (s *Storage) Register(p Provider) {
	s.Providers[p.Name()] = p
}

// RegisterAll registers all providers via the settings
func (s *Storage) RegisterAll(settings *models.Settings) {

	if settings.StorageOptions.Ftp.StorageOption.Active {
		s.Providers[settings.StorageOptions.Ftp.StorageOption.Option] = ftp.New(settings.StorageOptions.Ftp.Config)
	}

	if settings.StorageOptions.Dropbox.StorageOption.Active {
		s.Providers[settings.StorageOptions.Dropbox.StorageOption.Option] = dropbox.New(settings.StorageOptions.Dropbox.Config)
	}

	if settings.StorageOptions.GoogleDrive.StorageOption.Active {
		s.Providers[settings.StorageOptions.GoogleDrive.StorageOption.Option] = drive.New(settings.StorageOptions.GoogleDrive.Config)
	}

}

// ClearProviders will remove all providers currently in use.
func (s *Storage) ClearProviders() {
	s.Providers = make(map[string]Provider)
}

// GetProviders returns the registered providers
func (s *Storage) GetProviders() map[string]Provider {
	return s.Providers
}

// Reset resets the storage with the new settings
func (s *Storage) Reset(settings *models.Settings) error {
	s.ClearProviders()
	s.RegisterAll(settings)
	return nil
}

// GetProvider returns the registered providers
func (s *Storage) GetProvider(p string) (Provider, error) {
	provider, ok := s.Providers[p]
	if ok {
		return provider, nil
	}

	return nil, errors.Errorf("Provider %s not initialised", p)
}
