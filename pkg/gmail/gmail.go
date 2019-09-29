package gmail

import (
	"context"
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"io/ioutil"

	"github.com/mholt/archiver"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Gmail hold the client and oauth2 config
type Gmail struct {
	AuthConfig *oauth2.Config
	config     *config.Config
	db         *database.Store
	archiver   *archiver.Zip
}

// New init the Gmail struct and returns it
func New(cfg *config.Config, db *database.Store) (*Gmail, error) {
	b, err := ioutil.ReadFile(cfg.Google.File)
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	authCfg, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	g := &Gmail{}
	g.AuthConfig = authCfg
	g.config = cfg
	g.db = db
	g.archiver = archiver.NewZip()

	return g, nil
}

// GetAuthCodeURL returns the URL for getting the authorization code from the user
func (g *Gmail) GetAuthCodeURL() string {
	return g.AuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

// getClient gets the client for the given user and saves the token for the given user and saves it if it has been changed
func (g *Gmail) getClient(ac models.Account) (*gmail.Service, error) {

	if ac.OauthToken == nil && ac.GoogleToken != "" {
		t, err := g.AuthConfig.Exchange(context.TODO(), ac.GoogleToken)
		if err != nil {
			return nil, err
		}

		err = g.db.UpdateTokenAccount(&ac, t)
		if err != nil {
			return nil, err
		}
	}

	tokenSource := g.AuthConfig.TokenSource(oauth2.NoContext, ac.OauthToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if newToken.AccessToken != ac.OauthToken.AccessToken {
		err = g.db.UpdateTokenAccount(&ac, newToken)
		if err != nil {
			return nil, err
		}
	}

	client := oauth2.NewClient(oauth2.NoContext, tokenSource)

	api, err := gmail.New(client)
	if err != nil {
		return nil, errors.Errorf("Unable to retrieve Gmail client: %v", err)
	}

	return api, nil
}
