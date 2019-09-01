package gmail

import (
	"context"
	"gmail_backup/pkg/models"

	"golang.org/x/oauth2"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac *models.Account) (*oauth2.Token, error) {

	token, err := g.AuthConfig.Exchange(context.TODO(), ac.GoogleToken)
	if err != nil {
		return nil, err
	}

	return token, err

	// client := g.AuthConfig.Client(context.Background(), ac.GoogleToken.(*oauth2.Token))

	// return nil

}
