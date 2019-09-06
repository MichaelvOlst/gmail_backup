package gmail

import (
	"github.com/pkg/errors"
	"google.golang.org/api/gmail/v1"
)

// getUserLabels returns a list of the user labels
func (g *Gmail) getUserLabels(api *gmail.Service) ([]string, error) {

	user := "me"
	r, err := api.Users.Labels.List(user).Do()
	if err != nil {
		return nil, errors.Errorf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		return nil, errors.New("No labels found")
	}

	var labels []string
	for _, l := range r.Labels {
		labels = append(labels, l.Id)
	}

	return labels, nil
}
