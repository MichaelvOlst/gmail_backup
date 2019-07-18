package models

import (
	"time"
)

// Account ...
type Account struct {
	ID             int       `json:"id" storm:"id,increment"`
	Email          string    `json:"email" storm:"unique"`
	Attachments    bool      `json:"attachments"`
	BackupComplete bool      `json:"backup_complete"`
	BackupDate     time.Time `json:"backup_date"`
	AccessToken    string    `json:"accesstoken"`
}
