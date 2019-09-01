package gmail

import (
	"gmail_backup/pkg/config"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Gmail hold the client and oauth2 config
type Gmail struct {
	AuthConfig *oauth2.Config
	config     *config.Config
}

// New init the Gmail struct and returns it
func New(cfg *config.Config) (*Gmail, error) {
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

	return g, nil
}

// GetAuthCodeURL returns the URL for getting the authorization code from the user
func (g *Gmail) GetAuthCodeURL() string {
	return g.AuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
