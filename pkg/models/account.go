package models

import (
	"regexp"
	"time"

	"github.com/gobuffalo/validate"
)

// Account ...
type Account struct {
	ID             int       `json:"id" storm:"id,increment"`
	Email          string    `json:"email" storm:"unique"`
	EncryptionKey  string    `json:"encryption_key"`
	Attachments    bool      `json:"attachments"`
	BackupComplete bool      `json:"backup_complete"`
	BackupDate     time.Time `json:"backup_date"`
	AccessToken    string    `json:"accesstoken" `
}

// IsValid validates the struct attributes
func (a *Account) IsValid(errors *validate.Errors) {
	if a.EncryptionKey == "" {
		errors.Add("encryption_key", "Encyption key must not be blank")
	}

	if a.Email == "" {
		errors.Add("email", "Email must not be blank")
	}

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !rxEmail.MatchString(a.Email) {
		errors.Add("email", "Enter a valid email")
	}

	if a.AccessToken == "" {
		errors.Add("accesstoken", "Access token must not be blank")
	}
}
