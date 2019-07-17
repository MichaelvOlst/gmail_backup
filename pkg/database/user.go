package database

import "gmail_backup/pkg/models"

// CreateUser creates an user
func (s *Store) CreateUser(u *models.User) (*models.User, error) {
	err := s.Save(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteUser deletes an user
func (s *Store) DeleteUser(e string) error {
	var u models.User
	err := s.One("Email", e, &u)
	if err != nil {
		return err
	}

	err = s.DeleteStruct(&u)
	if err != nil {
		return err
	}

	return nil
}