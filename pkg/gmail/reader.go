package gmail

import (
	"fmt"
	"gmail_backup/pkg/database"
	"gmail_backup/pkg/models"
	"io"

	"github.com/dustin/go-humanize"
)

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
