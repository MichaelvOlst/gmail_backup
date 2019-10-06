package gmail

import (
	"archive/zip"
	"fmt"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"os"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
)

// Backup recieves an account to backup
func (g *Gmail) Backup(ac models.Account, s *storage.Storage) {

	if ac.BackupStarted == "true" {
		return
	}

	ac.BackupStarted = "true"
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

	// labels, err := g.getUserLabels(api)
	// if err != nil {
	// 	g.db.SaveAccountError(&ac, fmt.Sprintf("Could not get user labels: %v", err))
	// 	return
	// }

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
	lm, err := g.collectMessages(user, api, &ac)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Unable to collect messages: %v", err))
		return
	}

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

		if counter == 100 {
			break
		}

		_, err = g.db.SaveMessage(&models.Message{ID: m.Id, AccountID: ac.ID})
		if err != nil {
			g.db.SaveAccountResult(&ac, fmt.Sprintf("Could not save message with the id %s: %v", m.Id, err))
		}

		eta := 0.55 * float64((len(lm) - counter))
		g.db.SaveAccountResult(&ac, fmt.Sprintf("%d / %d %s", counter, totalMsg, secondsToHuman(int(eta))))

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

	err = g.upload(zipFilename, userPath, &ac, storage)
	if err != nil {
		g.db.SaveAccountError(&ac, fmt.Sprintf("Could not upload the zip: %v", err))
		return
	}

	g.db.SaveAccountResult(&ac, "Done")
	g.db.AccountBackupComplete(&ac)

}

func formatMessage(f string) googleapi.CallOption { return formatRequestMessage(f) }

type formatRequestMessage string

func (f formatRequestMessage) Get() (string, string) { return "format", string(f) }
