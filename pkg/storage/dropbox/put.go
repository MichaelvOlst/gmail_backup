package dropbox

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

// Put stores a file in Dropbox
func (p *Provider) Put(filename, path string, file *os.File, r io.Reader) error {

	contentsInfo, err := file.Stat()
	if err != nil {
		return err
	}

	path = strings.TrimRight(path, "/")
	filenameBase := fmt.Sprintf("%s/%s", path, filename)

	commitInfo := files.NewCommitInfo(filenameBase)
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	commitInfo.ClientModified = time.Now().UTC().Round(time.Second)

	dbx := p.client
	if contentsInfo.Size() > chunkSize {
		err = uploadChunked(dbx, r, commitInfo, contentsInfo.Size())
		if err != nil {
			return err
		}
		return nil
	}

	if _, err = dbx.Upload(commitInfo, file); err != nil {
		return err
	}

	return nil
}

func uploadChunked(dbx files.Client, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64) (err error) {
	res, err := dbx.UploadSessionStart(files.NewUploadSessionStartArg(), &io.LimitedReader{R: r, N: chunkSize})
	if err != nil {
		return
	}

	written := chunkSize

	for (sizeTotal - written) > chunkSize {
		cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
		args := files.NewUploadSessionAppendArg(cursor)

		err = dbx.UploadSessionAppendV2(args, &io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return
		}
		written += chunkSize
	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	if _, err = dbx.UploadSessionFinish(args, r); err != nil {
		return
	}

	return
}
