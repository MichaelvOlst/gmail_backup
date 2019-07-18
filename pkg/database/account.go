package database

import "gmail_backup/pkg/models"

// CreateAccount creates an account
func (s *Store) CreateAccount(a *models.Account) (*models.Account, error) {
	err := s.Save(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
