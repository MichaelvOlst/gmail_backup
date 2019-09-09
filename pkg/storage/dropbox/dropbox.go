package dropbox

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/ioprogress"
)

const name = "dropbox"
const chunkSize int64 = 1 << 24

// Config holds the config for the Dropbox option
type Config struct {
	AccessToken string `json:"accesstoken"`
}

// Provider implements storage.Provider for the dropbox file storage.
type Provider struct {
	client files.Client
}

// Name returns dropbox
func (p *Provider) Name() string {
	return name
}

// ListFolder returns dropbox
func (p *Provider) ListFolder() {
	arg := files.NewListFolderArg("")

	res, err := p.client.ListFolder(arg)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range res.Entries {
		switch f := entry.(type) {
		case *files.FileMetadata:
			fmt.Println(f.Name)
			// fmt.Printf("%+v\n", f)

			// printFileMetadata(w, f, long)
		case *files.FolderMetadata:
			fmt.Println(f.Name)
			// fmt.Printf("%+v\n", f)
			// printFolderMetadata(w, f, long)
		}
	}
}

// Put stores a file in Dropbox
func (p *Provider) Put(file string) {

	// fmt.Println("TODO " + file)
	contents, err := os.Open(file)
	if err != nil {
		return
	}
	defer contents.Close()

	contentsInfo, err := contents.Stat()
	if err != nil {
		return
	}

	progressbar := &ioprogress.Reader{
		Reader: contents,
		DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
			return fmt.Sprintf("Uploading %s/%s",
				humanize.IBytes(uint64(progress)), humanize.IBytes(uint64(total)))
		}),
		Size: contentsInfo.Size(),
	}

	// fmt.Println(filepath.Base(file))

	commitInfo := files.NewCommitInfo("/backup/" + filepath.Base(file))
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	commitInfo.ClientModified = time.Now().UTC().Round(time.Second)

	dbx := p.client
	if contentsInfo.Size() > chunkSize {
		err = uploadChunked(dbx, progressbar, commitInfo, contentsInfo.Size())
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err = dbx.Upload(commitInfo, progressbar); err != nil {
		log.Fatal(err)
		return
	}
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

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {
	config := dropbox.Config{
		Token: cfg.AccessToken,
	}

	c := files.New(config)

	return &Provider{c}
}
