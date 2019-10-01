package gmail

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"

	"google.golang.org/api/googleapi"

	"google.golang.org/api/gmail/v1"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac models.Account, s *storage.Storage) {

	if ac.BackupStarted != "done" || ac.BackupStarted == "" {
		return
	}

	ac.BackupStarted = "processing"
	err := g.db.Update(&ac)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not init backup process %v", err))
		return
	}

	// fmt.Printf("%+v\n", ac)

	dataFolder := "data"

	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		os.Mkdir(dataFolder, 0777)
	}

	api, err := g.getClient(ac)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not connect to gmail: %v", err))
		return
	}

	labels, err := g.getUserLabels(api)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not get user labels: %v", err))
		return
	}

	storage, err := s.GetProvider(ac.StorageProvider)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("%v", err))
		return
	}

	// err = g.db.Drop(&models.Message{})
	// if err != nil {
	// 	g.db.SaveAccountResult(&ac, fmt.Sprintf("Could not drop table messagess: %v", err))
	// 	return
	// }
	// fmt.Println("done dropping")
	// return

	userPath := fmt.Sprintf("/%s/%s", strings.TrimLeft(ac.UploadPath, "/"), ac.Email)
	if err = storage.Mkdir(userPath); err != nil && !storage.IsNotExists(err) {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Unable to create the folder: %s; error: %v", userPath, err))
		return
	}

	// var storedMsg []models.Message
	// err = g.db.All(&storedMsg)
	// if err != nil {
	// 	g.db.SaveAccountResult(&ac, fmt.Sprintf("Could not get list of messages; error: %v", err))
	// 	return
	// }
	// fmt.Printf("%+v", storedMsg)
	// return

	user := "me"
	lm := make(map[string]*gmail.Message)
	g.db.SaveAccountResult(&ac, "Collecting messages")

	for _, label := range labels {

		g.db.SaveAccountResult(&ac, fmt.Sprintf("Getting messages for label: %s", label))

		r, err := api.Users.Messages.List(user).IncludeSpamTrash(false).LabelIds(label).Do()
		if err != nil {
			g.db.SaveAccountError(&ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
			return
		}

		if len(r.Messages) == 0 {
			g.db.SaveAccountResult(&ac, fmt.Sprintf("No messages found with the label: %s", label))
			continue
		}

		for _, m := range r.Messages {

			sm, err := g.db.GetMessageByID(ac.ID, m.Id)
			if err != nil && err != storm.ErrNotFound {
				g.db.SaveAccountError(&ac, fmt.Sprintf("Could not get message with ID %s. %v", m.Id, err))
				return
			}

			if sm == nil {
				if _, ok := lm[m.Id]; !ok {
					lm[m.Id] = m
				}
			} else {
				fmt.Println("Found message with ID " + m.Id)
			}

		}

		if r.NextPageToken != "" {

			counter := 0
			nextPageToken := r.NextPageToken
			for {

				if counter == 0 {
					break
				}

				g.db.SaveAccountResult(&ac, fmt.Sprintf("Messages: %d", len(lm)))
				r, err := api.Users.Messages.List(user).LabelIds(label).IncludeSpamTrash(false).PageToken(nextPageToken).IncludeSpamTrash(false).Do()

				if err != nil {
					g.db.SaveAccountError(&ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
					return
				}

				if r.NextPageToken == "" || len(r.Messages) == 0 {
					break
				}

				nextPageToken = r.NextPageToken

				// lm = append(lm, r.Messages...)
				for _, m := range r.Messages {

					sm, err := g.db.GetMessageByID(ac.ID, m.Id)
					if err != nil && err != storm.ErrNotFound {
						g.db.SaveAccountError(&ac, fmt.Sprintf("Could not get message with ID %s. %v", m.Id, err))
						return
					}

					if sm == nil {
						if _, ok := lm[m.Id]; !ok {
							lm[m.Id] = m
						}
					} else {
						fmt.Println("Found message with ID " + m.Id)
					}
				}

				counter++
			}
		}
	}

	// // return false, nil
	// fmt.Println(len(lm))
	// return

	// fmt.Println("")

	// g.db.SaveAccountResult(&ac, "Done")
	// return

	if len(lm) == 0 {
		g.db.SaveAccountError(&ac, "No messages found")
		return
	}

	totalMsg := len(lm)
	g.db.SaveAccountResult(&ac, fmt.Sprintf("Total messages %d", totalMsg))

	zipFilename := fmt.Sprintf("%s/%s-%s.zip", dataFolder, ac.Email, time.Now().Format("2006-01-02-15:04"))
	zipfile, err := os.Create(zipFilename)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not create temp %s file: %v", zipFilename, err))
		return
	}

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	counter := 0
	for _, m := range lm {
		md, err := api.Users.Messages.Get(user, m.Id).Do(formatMessage("raw"))
		if err != nil {
			g.db.SaveAccountResult(&ac, fmt.Sprintf("Could not retrieve message with the id %s: %v", m.Id, err))
		}

		err = g.saveMessage(userPath, md, zipWriter)
		if err != nil {
			g.db.SaveAccountResult(&ac, fmt.Sprintf("Could not store message with the id %s: %v", m.Id, err))
		}

		if counter == 10 {
			break
		}

		_, err = g.db.SaveMessage(&models.Message{ID: m.Id, AccountID: ac.ID})
		if err != nil {
			g.db.SaveAccountResult(&ac, fmt.Sprintf("Could not save message with the id %s: %v", m.Id, err))
		}

		eta := 0.55 * float64((len(lm) - counter))
		g.db.SaveAccountResult(&ac, fmt.Sprintf("%d 	/ %d %s", counter, totalMsg, secondsToHuman(int(eta))))

		counter++
	}

	lm = nil

	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not close the zip file: %v", err))
		return
	}
	zipfile.Close()

	zipfile, err = os.Open(zipFilename)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not open temp zip file: %v", err))
		return
	}
	defer zipfile.Close()

	zipFileStats, err := zipfile.Stat()
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not get stat from zip file: %v", err))
		return
	}

	r := &progressReader{
		Reader:    zipfile,
		TotalSize: zipFileStats.Size(),
		db:        g.db,
		account:   &ac,
	}

	dropboxName := fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02-15:04"))

	storage.Put(dropboxName, userPath, zipfile, r)

	err = os.Remove(zipFilename)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not remove temp zip file %s: %v", zipFilename, err))
		return
	}

	g.db.SaveAccountResult(&ac, "Done")
	g.db.AccountBackupComplete(&ac)

	// ac.BackupStarted = 0
	// err = g.db.Update(&ac)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// g.db.SaveAccountError(&ac, fmt.Sprintf("Could not init backup process %v", err))
	// 	// return
	// }

}

func (g *Gmail) saveMessage(path string, m *gmail.Message, zipWriter *zip.Writer) error {

	b, _ := g.decodeMessage(m)
	t := time.Unix(m.InternalDate/1000, 0)
	folder := fmt.Sprintf("%s", t.Format("2006-01"))

	tempFile, err := ioutil.TempFile("", "backup_")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(b)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s.eml", folder, m.Id)

	info, err := tempFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	writer.Write(b)
	_, err = io.Copy(writer, tempFile)
	return err
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
