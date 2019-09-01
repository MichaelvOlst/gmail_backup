package database

import (
	"gmail_backup/pkg/models"

	"golang.org/x/oauth2"
)

// CreateAccount creates an account
func (s *Store) CreateAccount(a *models.Account) (*models.Account, error) {
	err := s.Save(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAccountByID gets an account by id
func (s *Store) GetAccountByID(id int) (*models.Account, error) {
	var ac models.Account
	err := s.One("ID", id, &ac)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

// UpdateAccount updates an account
func (s *Store) UpdateAccount(ac *models.Account) (*models.Account, error) {
	err := s.Update(ac)
	if err != nil {
		return nil, err
	}
	return ac, nil
}

// DeleteAccount updates an account
func (s *Store) DeleteAccount(ac *models.Account) error {
	err := s.DeleteStruct(ac)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTokenAccount updates the token for the given account
func (s *Store) UpdateTokenAccount(ac *models.Account, t *oauth2.Token) error {
	ac.OauthToken = t
	err := s.Update(ac)
	if err != nil {
		return err
	}
	return nil
}
