package account

import (
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"

	"github.com/robfig/cron/v3"
)

// Account manages the cronjob and db etc..
type Account struct {
	db      *database.Store
	cronjob *cron.Cron
	storage *storage.Storage
	config  *config.Config
}

// New returns an instance of Account
func New(db *database.Store, storage *storage.Storage, config *config.Config) *Account {
	return &Account{
		db:      db,
		storage: storage,
		config:  config,
		cronjob: cron.New(),
	}
}

func (a *Account) getAllAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := a.db.All(&accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
