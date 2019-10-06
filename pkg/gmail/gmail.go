package gmail

import (
	"context"
	"fmt"
	"gmail_backup/pkg/config"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"io/ioutil"
	"math"

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

func secondsToHuman(input int) string {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	var result string
	if years > 0 {
		result = formatTime(int(years), "y") + formatTime(int(months), "m") + formatTime(int(weeks), "w") + formatTime(int(days), "d") + formatTime(int(hours), "h") + formatTime(int(minutes), "m") + formatTime(int(seconds), "s")
	} else if months > 0 {
		result = formatTime(int(months), "m") + formatTime(int(weeks), "w") + formatTime(int(days), "d") + formatTime(int(hours), "h") + formatTime(int(minutes), "m") + formatTime(int(seconds), "s")
	} else if weeks > 0 {
		result = formatTime(int(weeks), "w") + formatTime(int(days), "d") + formatTime(int(hours), "h") + formatTime(int(minutes), "m") + formatTime(int(seconds), "s")
	} else if days > 0 {
		result = formatTime(int(days), "d") + formatTime(int(hours), "h") + formatTime(int(minutes), "m") + formatTime(int(seconds), "s")
	} else if hours > 0 {
		result = formatTime(int(hours), "h") + formatTime(int(minutes), "m") + formatTime(int(seconds), "s")
	} else if minutes > 0 {
		result = formatTime(int(minutes), "m") + formatTime(int(seconds), "s")
	} else {
		result = formatTime(int(seconds), "s")
	}

	return result
}

func formatTime(count int, format string) string {
	return fmt.Sprintf("%d%s", count, format)
}
