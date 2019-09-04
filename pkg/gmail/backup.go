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

	var lm []*gmail.Message

	r, err := api.Users.Messages.List(user).Do()
	if err != nil {
		return false, errors.Errorf("Unable to retrieve messages: %v", err)
	}

	if len(r.Messages) == 0 {
		return false, errors.New("No messages found")
	}

	lm = append(lm, r.Messages...)

	if r.NextPageToken != "" {

		fmt.Println("next set")

		counter := 0
		nextPageToken := r.NextPageToken
		for {

			if counter == 100 {
				break
			}

			fmt.Printf("Getting next messages %d\n", counter)
			r, err := api.Users.Messages.List(user).Do(pageToken(nextPageToken))

			if err != nil {
				return false, errors.Errorf("Unable to retrieve messages: %v", err)
			}

			if r.NextPageToken == "" || len(r.Messages) == 0 {
				break
			}

			nextPageToken = r.NextPageToken

			lm = append(lm, r.Messages...)

			counter++
		}
	}

	// return false, nil
	// fmt.Println(len(lm))
	// fmt.Println("")

	if len(lm) == 0 {
		return false, errors.New("No messages found")
	}

	ch := make(chan message, len(lm))
	defer close(ch)

	for _, m := range lm {
		fmt.Printf("Getting message with id: %s\n", m.Id)
		go func(api *gmail.Service, user string, m *gmail.Message, ch chan<- message) {

			var err error
			md, err := api.Users.Messages.Get(user, m.Id).Do(formatMessage("raw"))
			if err != nil {
				ch <- message{
					id:      m.Id,
					snippet: m.Snippet,
					err:     err,
				}
			} else {
				ch <- message{
					id:           m.Id,
					internalDate: md.InternalDate,
					snippet:      md.Snippet,
					raw:          md.Raw,
				}
			}

		}(api, user, m, ch)
	}

	for m := range ch {
		if m.err != nil {
			fmt.Printf("SHow me the error %v\n", m.err)
		} else {
			err = g.saveMessage(path, m)
			if err != nil {
				return false, err
			}
		}

	}

	fmt.Println("Done")

	// for _, l := range r.Messages {
	// 	go saveMessage()

	// 	m, err := api.Users.Messages.Get(user, l.Id).Do(formatMessage("raw"))
	// 	if err != nil {
	// 		return false, errors.Errorf("Unable to retrieve message: %v", err)
	// 	}

	// 	err = g.saveMessage(path, m)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	break
	// }

	return true, nil
}

type message struct {
	id           string
	internalDate int64
	snippet      string
	raw          string
	err          error
}

func (g *Gmail) saveMessage(path string, m message) error {

	b, _ := g.decodeMessage(m)

	t := time.Unix(m.internalDate/1000, 0)

	path = fmt.Sprintf("%s/%s", path, t.Format("2006-01"))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	filename := fmt.Sprintf("%s/%s.eml", path, m.id)
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

func (g *Gmail) decodeMessage(m message) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(m.raw)
}

func formatMessage(f string) googleapi.CallOption { return formatRequestMessage(f) }

type formatRequestMessage string

func (f formatRequestMessage) Get() (string, string) { return "format", string(f) }

func pageToken(f string) googleapi.CallOption { return pageTokenOption(f) }

type pageTokenOption string

func (f pageTokenOption) Get() (string, string) { return "pageToken", string(f) }
