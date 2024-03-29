package models

import (
	"regexp"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// Account is the model for accounts
type Account struct {
	ID                int           `json:"id" storm:"id,increment"`
	Email             string        `json:"email" storm:"unique"`
	EncryptionKey     string        `json:"encryption_key"`
	BackupComplete    bool          `json:"backup_complete"`
	BackupStarted     string        `json:"backup_started"`
	BackupDate        time.Time     `json:"backup_date"`
	BackupProgressMsg string        `json:"backup_progress_message"`
	GoogleToken       string        `json:"google_token"`
	OauthToken        *oauth2.Token `json:"token"`
	StorageProvider   string        `json:"storage_provider"`
	UploadPath        string        `json:"upload_path"`
	CronExpression    string        `json:"cron_expression"`
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

	if strings.TrimSpace(a.StorageProvider) == "" {
		v["storage_provider"] = "Storage option is required"
	}

	if strings.TrimSpace(a.UploadPath) == "" {
		v["upload_path"] = "Upload path is required"
	}

	if strings.TrimSpace(a.CronExpression) == "" {
		v["cron_expression"] = "Cronjob is required"
	}

	return v
}
