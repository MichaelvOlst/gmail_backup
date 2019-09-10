package gmail

import (
	"encoding/base64"
	"fmt"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/gmail/v1"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac *models.Account, s *storage.Storage) {

	// api, err := g.getClient(ac)
	// if err != nil {
	// 	g.db.SaveAccountResult(ac, fmt.Sprintf("Could not connect to gmail: %v", err))
	// 	return
	// }

	// labels, err := g.getUserLabels(api)
	// if err != nil {
	// 	g.db.SaveAccountResult(ac, fmt.Sprintf("Could not get user labels: %v", err))
	// 	return
	// }

	storage, err := s.GetProvider(ac.StorageProvider)
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("%v", err))
		return
	}

	userPath := fmt.Sprintf("%s/%s", strings.TrimLeft(g.config.Backup.Path, "/"), ac.Email)
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		os.Mkdir(userPath, 0777)
	}

	// user := "me"

	// // var lm []*gmail.Message
	// lm := make(map[string]*gmail.Message)

	// g.db.SaveAccountResult(ac, "Collecting messages")

	// for _, label := range labels {

	// 	g.db.SaveAccountResult(ac, fmt.Sprintf("Getting messages for label: %s", label))

	// 	r, err := api.Users.Messages.List(user).IncludeSpamTrash(false).LabelIds(label).Do()
	// 	if err != nil {
	// 		g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
	// 		return
	// 	}

	// 	if len(r.Messages) == 0 {
	// 		g.db.SaveAccountResult(ac, fmt.Sprintf("No messages found with the label: %s", label))
	// 		continue
	// 	}

	// 	for _, m := range r.Messages {
	// 		if _, ok := lm[m.Id]; !ok {
	// 			lm[m.Id] = m
	// 		}
	// 	}

	// 	// lm = append(lm, r.Messages...)

	// 	if r.NextPageToken != "" {

	// 		counter := 0
	// 		nextPageToken := r.NextPageToken
	// 		for {

	// 			// if counter == 0 {
	// 			// 	break
	// 			// }

	// 			g.db.SaveAccountResult(ac, fmt.Sprintf("Messages: %d", len(lm)))
	// 			r, err := api.Users.Messages.List(user).LabelIds(label).IncludeSpamTrash(false).PageToken(nextPageToken).IncludeSpamTrash(false).Do()

	// 			if err != nil {
	// 				g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
	// 				return
	// 			}

	// 			if r.NextPageToken == "" || len(r.Messages) == 0 {
	// 				break
	// 			}

	// 			nextPageToken = r.NextPageToken

	// 			// lm = append(lm, r.Messages...)
	// 			for _, m := range r.Messages {
	// 				if _, ok := lm[m.Id]; !ok {
	// 					lm[m.Id] = m
	// 				}
	// 			}

	// 			counter++
	// 		}
	// 	}
	// }

	// // return false, nil
	// fmt.Println(len(lm))
	// return

	// fmt.Println("")

	// g.db.SaveAccountResult(ac, "Done")
	// return

	// if len(lm) == 0 {
	// 	g.db.SaveAccountResult(ac, "No messages found")
	// 	return
	// }

	// totalMsg := len(lm)
	// g.db.SaveAccountResult(ac, fmt.Sprintf("Total messages %d", totalMsg))

	// counter := 0
	// for _, m := range lm {

	// 	md, err := api.Users.Messages.Get(user, m.Id).Do(formatMessage("raw"))
	// 	if err != nil {
	// 		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not retrieve message with the id %s: %v", m.Id, err))
	// 	}

	// 	err = g.saveMessage(userPath, md)
	// 	if err != nil {
	// 		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not save message with the id %s: %v", m.Id, err))
	// 	}

	// 	eta := 0.55 * float64((len(lm) - counter))
	// 	g.db.SaveAccountResult(ac, fmt.Sprintf("%d 	/ %d %s", counter, totalMsg, secondsToHuman(int(eta))))

	// 	counter++
	// }

	// lm = make(map[string]*gmail.Message) // making map empty

	g.db.SaveAccountResult(ac, "Zipping messages")

	root := userPath
	t := time.Now()
	zipPath := fmt.Sprintf("%s/%s-%s.zip", userPath, ac.Email, t.Format("2006-01-02"))

	if _, err := os.Stat(zipPath); err == nil {
		err = os.Remove(zipPath)
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Couldn't remove existing zip file: %v", err))
		}
	}

	g.db.SaveAccountResult(ac, "Zipping messages..")
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return g.archiver.Archive([]string{path}, zipPath)
	})
	// if err != nil {
	// 	g.db.SaveAccountResult(ac, fmt.Sprintf("Error adding file to archive: %v", err))
	// 	return
	// }

	g.db.SaveAccountResult(ac, fmt.Sprintf("Saving %s to %s", filepath.Base(zipPath), storage.Name()))

	contents, err := os.Open(zipPath)
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not open file: %v", err))
		return
	}
	defer contents.Close()

	contentsInfo, err := contents.Stat()
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not get file info: %v", err))
		return
	}

	r := &progressReader{
		Reader:    contents,
		TotalSize: contentsInfo.Size(),
		db:        g.db,
		account:   ac,
	}

	// progressbar := &ioprogress.Reader{
	// 	Reader: contents,
	// 	DrawFunc: ioprogress.DrawTerminalf(os.Stdout, func(progress, total int64) string {
	// 		ur := fmt.Sprintf("Uploading %s/%s", humanize.Bytes(uint64(progress)), humanize.Bytes(uint64(total)))
	// 		return ur
	// 	}),
	// 	Size: contentsInfo.Size(),
	// }

	storage.Put(zipPath, ac.UploadPath, r)

	g.db.SaveAccountResult(ac, "Done")

}

type progressReader struct {
	io.Reader
	written   int64 // Total # of bytes transferred
	TotalSize int64
	db        *database.Store
	account   *models.Account
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if err == nil {
		pr.written += int64(n)
		ur := fmt.Sprintf("Uploading %s/%s", humanize.Bytes(uint64(pr.written)), humanize.Bytes(uint64(pr.TotalSize)))
		pr.db.SaveAccountResult(pr.account, ur)
	}
	return n, err
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
