package dropbox

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
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
func (p *Provider) Put(filename, path string, file *os.File) error {

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
		err = uploadChunked(dbx, file, commitInfo, contentsInfo.Size())
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

// Mkdir creates a folder in dropbox
func (p *Provider) Mkdir(path string) error {
	arg := files.NewCreateFolderArg(path)

	_, err := p.client.CreateFolderV2(arg)
	if err != nil {
		return err
	}
	return nil
}

// IsNotExists check if a folder already exists
func (p *Provider) IsNotExists(err error) bool {
	cerr, ok := err.(files.CreateFolderAPIError)
	if !ok {
		return false
	}

	if strings.Contains(cerr.APIError.Error(), "path/conflict/folder/") {
		return true
	}

	fmt.Println(cerr.APIError)

	return false
}

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {
	config := dropbox.Config{
		Token: cfg.AccessToken,
	}

	c := files.New(config)

	return &Provider{c}
}
