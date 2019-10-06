package gmail

import (
	"fmt"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"gmail_backup/pkg/storage"
	"io"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

func (g *Gmail) upload(file, userPath string, ac *models.Account, storage storage.Provider) error {
	zipfile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	zipFileStats, err := zipfile.Stat()
	if err != nil {
		return err
	}

	r := &progressReader{
		Reader:    zipfile,
		TotalSize: zipFileStats.Size(),
		db:        g.db,
		account:   ac,
	}

	dropboxName := fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02-15:04"))

	storage.Put(dropboxName, userPath, zipfile, r)

	err = os.Remove(file)
	if err != nil {
		return err
	}

	return nil
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
