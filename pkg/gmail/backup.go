package gmail

import (
	"compress/flate"
	"encoding/base64"
	"fmt"
	"gmail_backup/pkg/models"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"

	"github.com/mholt/archiver"
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

	userPath := fmt.Sprintf("%s/%s", strings.TrimLeft(g.config.Backup.Path, "/"), ac.Email)
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		os.Mkdir(userPath, 0777)
	}

	user := "me"

	// var lm []*gmail.Message
	lm := make(map[string]*gmail.Message)

	g.db.SaveAccountResult(ac, "Getting messages...")

	_ = labels
	// for _, label := range labels {

	label := "INBOX"

	g.db.SaveAccountResult(ac, fmt.Sprintf("Getting messages for label: %s", label))

	r, err := api.Users.Messages.List(user).IncludeSpamTrash(false).LabelIds(label).Do()
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
		return
	}

	if len(r.Messages) == 0 {
		g.db.SaveAccountResult(ac, fmt.Sprintf("No messages found with the label: %s", label))
		// continue
	}

	for _, m := range r.Messages {
		if _, ok := lm[m.Id]; !ok {
			lm[m.Id] = m
		}
	}

	// lm = append(lm, r.Messages...)

	// if r.NextPageToken != "" {

	// 	counter := 0
	// 	nextPageToken := r.NextPageToken
	// 	for {

	// 		if counter == 0 {
	// 			break
	// 		}

	// 		g.db.SaveAccountResult(ac, fmt.Sprintf("Total messages: %d", len(lm)))
	// 		r, err := api.Users.Messages.List(user).LabelIds(label).IncludeSpamTrash(false).PageToken(nextPageToken).IncludeSpamTrash(false).Do()

	// 		if err != nil {
	// 			g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
	// 			return
	// 		}

	// 		if r.NextPageToken == "" || len(r.Messages) == 0 {
	// 			break
	// 		}

	// 		nextPageToken = r.NextPageToken

	// 		// lm = append(lm, r.Messages...)
	// 		for _, m := range r.Messages {
	// 			if _, ok := lm[m.Id]; !ok {
	// 				lm[m.Id] = m
	// 			}
	// 		}

	// 		counter++
	// 	}
	// }
	// }

	// // return false, nil
	// fmt.Println(len(lm))
	// return

	// fmt.Println("")

	// g.db.SaveAccountResult(ac, "Done")
	// return

	if len(lm) == 0 {
		g.db.SaveAccountResult(ac, "No messages found")
		return
	}

	totalMsg := len(lm)
	g.db.SaveAccountResult(ac, fmt.Sprintf("Total messages %d", totalMsg))

	counter := 0
	for _, m := range lm {

		md, err := api.Users.Messages.Get(user, m.Id).Do(formatMessage("raw"))
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Could not retrieve message with the id %s: %v", m.Id, err))
		}

		err = g.saveMessage(userPath, md)
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Could not save message with the id %s: %v", m.Id, err))
		}

		eta := 0.55 * float64((len(lm) - counter))
		g.db.SaveAccountResult(ac, fmt.Sprintf("%d 	/ %d %s", counter, totalMsg, secondsToHuman(int(eta))))

		counter++
	}

	lm = make(map[string]*gmail.Message) // making map empty

	g.db.SaveAccountResult(ac, "Zipping messages")

	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}
	root := userPath
	zipPath := userPath + "/test.zip"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return z.Archive([]string{path}, zipPath)
		}
		return nil
	})
}

func (g *Gmail) saveMessage(path string, m *gmail.Message) error {

	b, _ := g.decodeMessage(m)

	t := time.Unix(m.InternalDate/1000, 0)

	path = fmt.Sprintf("%s/%s", path, t.Format("2006-01"))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	filename := fmt.Sprintf("%s/%s.eml", path, m.Id)
	// fmt.Println(filename)
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

func formatMessage(f string) googleapi.CallOption { return formatRequestMessage(f) }

type formatRequestMessage string

func (f formatRequestMessage) Get() (string, string) { return "format", string(f) }
