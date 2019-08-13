package models

import (
	"regexp"
	"strings"
	"time"
)

// Account is the model for accounts
type Account struct {
	ID             int       `json:"id" storm:"id,increment"`
	Email          string    `json:"email" storm:"unique"`
	EncryptionKey  string    `json:"encryption_key"`
	Attachments    bool      `json:"attachments"`
	BackupComplete bool      `json:"backup_complete"`
	BackupDate     time.Time `json:"backup_date"`
	GoogleToken    string    `json:"google_token"`
}

// Validate the account model
func (a *Account) Validate() map[string]string {
	var v = make(map[string]string)

	if strings.TrimSpace(a.EncryptionKey) == "" {
		v["encryption_key"] = "Encyption key is required"
	}

	if strings.TrimSpace(a.GoogleToken) == "" {
		v["google_token"] = "Google token is required"
	}

	if strings.TrimSpace(a.Email) == "" {
		v["email"] = "Email is required"
	} else {
		var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !rxEmail.MatchString(a.Email) {
			v["email"] = "Enter an valid email"
		}
	}

	return v
}
