package database

import (
	"gmail_backup/pkg/models"

	"github.com/asdine/storm/q"
)

// ListMessages saves a new message
func (s *Store) ListMessages(m *models.Message) (*models.Message, error) {
	err := s.Save(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// SaveMessage saves a new message
func (s *Store) SaveMessage(m *models.Message) (*models.Message, error) {
	err := s.Save(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetMessageByID gets a single message by ID
func (s *Store) GetMessageByID(accountID int, ID string) (*models.Message, error) {
	var m models.Message
	err := s.Select(q.Eq("AccountID", accountID), q.Eq("ID", ID)).First(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
