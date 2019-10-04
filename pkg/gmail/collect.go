package gmail

import (
	"fmt"
	"gmail_backup/pkg/models"

	"github.com/asdine/storm"
	"github.com/pkg/errors"
	"google.golang.org/api/gmail/v1"
)

func (g *Gmail) collectMessages(user string, api *gmail.Service, ac *models.Account) (map[string]*gmail.Message, error) {

	labels, err := g.getUserLabels(api)
	if err != nil {
		return nil, errors.Errorf("Could not get user labels: %v", err)
	}

	lm := make(map[string]*gmail.Message)

	for _, label := range labels {

		r, err := api.Users.Messages.List(user).IncludeSpamTrash(false).LabelIds(label).Do()
		if err != nil {
			return nil, errors.Errorf("Unable to retrieve messages: %v", err)
		}

		if len(r.Messages) == 0 {
			continue
		}

		for _, m := range r.Messages {

			sm, err := g.db.GetMessageByID(ac.ID, m.Id)
			if err != nil && err != storm.ErrNotFound {
				return nil, errors.Errorf("Error getting message %s from db: %v", m.Id, err)
			}

			if sm == nil {
				if _, ok := lm[m.Id]; !ok {
					lm[m.Id] = m
				}
			}
		}

		if r.NextPageToken != "" {

			counter := 0
			nextPageToken := r.NextPageToken
			for {

				// if counter == 0 {
				// 	break
				// }

				g.db.SaveAccountResult(ac, fmt.Sprintf("Messages: %d", len(lm)))
				r, err := api.Users.Messages.List(user).LabelIds(label).IncludeSpamTrash(false).PageToken(nextPageToken).IncludeSpamTrash(false).Do()

				if err != nil {
					return nil, errors.Errorf("Unable to retrieve messages: %v", err)
				}

				if r.NextPageToken == "" || len(r.Messages) == 0 {
					break
				}

				nextPageToken = r.NextPageToken

				for _, m := range r.Messages {

					sm, err := g.db.GetMessageByID(ac.ID, m.Id)
					if err != nil && err != storm.ErrNotFound {
						return nil, errors.Errorf("Error getting message %s from db: %v", m.Id, err)
					}

					if sm == nil {
						if _, ok := lm[m.Id]; !ok {
							lm[m.Id] = m
						}
					}
				}

				counter++
			}
		}
	}
	return lm, nil
}
