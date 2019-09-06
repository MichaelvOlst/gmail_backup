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
func (g *Gmail) Backup(ac *models.Account) {

	api, err := g.getClient(ac)
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not connect to gmail: %v", err))
		return
	}

	labels, err := g.getUserLabels(api)
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not get user labels: %v", err))
		return
	}

	path := fmt.Sprintf("%s/%s", strings.TrimLeft(g.config.Backup.Path, "/"), ac.Email)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	user := "me"

	var lm []*gmail.Message

	g.db.SaveAccountResult(ac, "Getting messages...")

	for _, label := range labels {

		g.db.SaveAccountResult(ac, fmt.Sprintf("Getting messages for label: %s", label))

		r, err := api.Users.Messages.List(user).LabelIds(label).Do()
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
			return
		}

		if len(r.Messages) == 0 {
			g.db.SaveAccountResult(ac, fmt.Sprintf("No messages found with the label: %s", label))
			continue
		}

		lm = append(lm, r.Messages...)

		if r.NextPageToken != "" {

			counter := 0
			nextPageToken := r.NextPageToken
			for {

				// if counter == 1 {
				// 	break
				// }

				g.db.SaveAccountResult(ac, fmt.Sprintf("Getting messages %d", len(lm)))
				r, err := api.Users.Messages.List(user).PageToken(nextPageToken).IncludeSpamTrash(false).Do()

				if err != nil {
					g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
					return
				}

				if r.NextPageToken == "" || len(r.Messages) == 0 {
					break
				}

				nextPageToken = r.NextPageToken

				lm = append(lm, r.Messages...)

				counter++
			}
		}
	}

	// // return false, nil
	// fmt.Println(len(lm))
	// return

	// fmt.Println("")

	if len(lm) == 0 {
		g.db.SaveAccountResult(ac, "No messages found")
		return
	}

	totalMsg := len(lm)
	g.db.SaveAccountResult(ac, fmt.Sprintf("Total messages %d", totalMsg))

	for k, m := range lm {
		g.db.SaveAccountResult(ac, fmt.Sprintf("%d / %d", k, totalMsg))

		md, err := api.Users.Messages.Get(user, m.Id).Do(formatMessage("raw"))
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Could not retrieve message with the id %s: %v", m.Id, err))
		}

		err = g.saveMessage(path, md)
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Could not save message with the id %s: %v", m.Id, err))
		}
	}

	g.db.SaveAccountResult(ac, "Done")
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
