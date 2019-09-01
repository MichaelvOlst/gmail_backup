package gmail

import (
	"fmt"
	"gmail_backup/pkg/models"
	"log"

	"google.golang.org/api/googleapi"

	"google.golang.org/api/gmail/v1"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac *models.Account) error {

	client, err := g.getClient(ac)
	if err != nil {
		return err
	}
	_ = client

	api, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"

	r, err := api.Users.Messages.List(user).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}
	if len(r.Messages) == 0 {
		fmt.Println("No labels found.")
		return nil
	}

	for _, l := range r.Messages {
		m, err := api.Users.Messages.Get(user, l.Id).Do(format("raw"))
		if err != nil {
			log.Fatalf("Unable to retrieve message: %v", err)
		}

		fmt.Println(m.Raw)
		break
	}

	return nil

}

func format(f string) googleapi.CallOption { return formatMessage(f) }

type formatMessage string

func (f formatMessage) Get() (string, string) { return "format", string(f) }
