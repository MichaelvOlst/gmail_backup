package gmail

import (
	"encoding/base64"
	"fmt"
	"gmail_backup/pkg/models"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/gmail/v1"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac *models.Account) (bool, error) {

	client, err := g.getClient(ac)
	if err != nil {
		return false, err
	}

	api, err := gmail.New(client)
	if err != nil {
		return false, errors.Errorf("Unable to retrieve Gmail client: %v", err)
	}

	path := fmt.Sprintf("%s/%s", strings.TrimLeft(g.config.Backup.Path, "/"), ac.Email)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	user := "me"
	r, err := api.Users.Messages.List(user).Do()
	if err != nil {
		return false, errors.Errorf("Unable to retrieve messages: %v", err)
	}
	if len(r.Messages) == 0 {
		return false, errors.New("No messages found")
	}

	for _, l := range r.Messages {
		m, err := api.Users.Messages.Get(user, l.Id).Do(formatMessage("raw"))
		if err != nil {
			return false, errors.Errorf("Unable to retrieve message: %v", err)
		}

		err = g.saveMessage(path, m)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (g *Gmail) saveMessage(path string, m *gmail.Message) error {

	b, _ := g.decodeMessage(m)

	t := time.Unix(m.InternalDate/1000, 0)

	path = fmt.Sprintf("%s/%s", path, t.Format("2006-01"))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	filename := fmt.Sprintf("%s/%s.eml", path, m.Id)
	fmt.Println(filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return errors.Errorf("Unable to create raw message file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gmail) decodeMessage(m *gmail.Message) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(m.Raw)
}

func formatMessage(f string) googleapi.CallOption { return formatRequestMessage(f) }

type formatRequestMessage string

func (f formatRequestMessage) Get() (string, string) { return "format", string(f) }
