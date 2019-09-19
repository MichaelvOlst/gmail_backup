package gmail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"github.com/mholt/archiver"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/gmail/v1"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac *models.Account, s *storage.Storage) {

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

	storage, err := s.GetProvider(ac.StorageProvider)
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("%v", err))
		return
	}

	if ac.UploadPath == "" {
		ac.UploadPath = "/test/"
	}

	userPath := fmt.Sprintf("/%s/%s", strings.TrimLeft(ac.UploadPath, "/"), ac.Email)
	if err = storage.Mkdir(userPath); err != nil && !storage.IsNotExists(err) {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to create the folder: %s; error: %v", userPath, err))
		return
	}

	user := "me"
	lm := make(map[string]*gmail.Message)
	g.db.SaveAccountResult(ac, "Collecting messages")

	for _, label := range labels {

		g.db.SaveAccountResult(ac, fmt.Sprintf("Getting messages for label: %s", label))

		r, err := api.Users.Messages.List(user).IncludeSpamTrash(false).LabelIds(label).Do()
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Unable to retrieve messages: %v", err))
			return
		}

		if len(r.Messages) == 0 {
			g.db.SaveAccountResult(ac, fmt.Sprintf("No messages found with the label: %s", label))
			continue
		}

		for _, m := range r.Messages {
			if _, ok := lm[m.Id]; !ok {
				lm[m.Id] = m
			}
		}

		break

		// if len(lm) > 100 {
		// 	break
		// }

		// lm = append(lm, r.Messages...)

		// if r.NextPageToken != "" {

		// 	counter := 0
		// 	nextPageToken := r.NextPageToken
		// 	for {

		// 		if counter == 0 {
		// 			break
		// 		}

		// 		g.db.SaveAccountResult(ac, fmt.Sprintf("Messages: %d", len(lm)))
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
	}

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

	buf := new(bytes.Buffer)
	zipWriter := bufio.NewWriter(buf)

	z := archiver.NewZip()
	err = z.Create(zipWriter)
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not create zip file: %v", err))
		return
	}
	defer z.Close()

	counter := 0
	for _, m := range lm {
		md, err := api.Users.Messages.Get(user, m.Id).Do(formatMessage("raw"))
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Could not retrieve message with the id %s: %v", m.Id, err))
		}

		err = g.saveMessage(userPath, md, z)
		if err != nil {
			g.db.SaveAccountResult(ac, fmt.Sprintf("Could not save message with the id %s: %v", m.Id, err))
		}

		if counter == 10 {
			break
		}

		eta := 0.55 * float64((len(lm) - counter))
		g.db.SaveAccountResult(ac, fmt.Sprintf("%d 	/ %d %s", counter, totalMsg, secondsToHuman(int(eta))))

		counter++
	}

	lm = make(map[string]*gmail.Message) // making map empty

	// fmt.Printf("%v\n", zipWriter)
	file, err := ioutil.TempFile("", "zip_")
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not create temp zip file: %v", err))
		return
	}
	defer file.Close()

	// io.Copy(zipWriter, file)

	b, err := file.Write(buf.Bytes())
	if err != nil {
		g.db.SaveAccountResult(ac, fmt.Sprintf("Could not write to temp zip file: %v", err))
		return
	}
	fmt.Println(b)
	// file.Write(zipWriter)
	// fmt.Printf("%v\n", byteBuffer.Bytes())

	storage.Put("/test.zip", ac.UploadPath, file)

	g.db.AccountBackupComplete(ac)

	g.db.SaveAccountResult(ac, "Done")
}

func (g *Gmail) saveMessage(path string, m *gmail.Message, z *archiver.Zip) error {

	b, _ := g.decodeMessage(m)
	t := time.Unix(m.InternalDate/1000, 0)
	path = fmt.Sprintf("%s/%s", path, t.Format("2006-01"))

	file, err := ioutil.TempFile("", "backup_")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	info, err := os.Stat(file.Name())
	if err != nil {
		return err
	}

	// get file's name for the inside of the archive
	internalName, err := archiver.NameInArchive(info, m.Id+".eml", path)
	if err != nil {
		return err
	}

	// write it to the archive
	err = z.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   info,
			CustomName: internalName,
		},
		ReadCloser: file,
	})

	if err != nil {
		return err
	}

	// if err := storage.Mkdir(path); !storage.IsNotExists(err) {
	// 	return errors.Errorf("Unable to create folder: %s", path)
	// }

	// filename := fmt.Sprintf("/%s.eml", m.Id)

	// // fmt.Println(filename)

	// file, err := ioutil.TempFile("", "backup_")
	// if err != nil {
	// 	return err
	// }
	// defer os.Remove(file.Name())

	// _, err = file.Write(b)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(path)

	// storage.Put(filename, path, file)
	// fmt.Println(filename)
	// f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	// if err != nil {
	// 	return errors.Errorf("Unable to create raw message file: %v", err)
	// }
	// defer f.Close()

	// _, err = f.Write(b)
	// if err != nil {
	// 	return err
	// }

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
