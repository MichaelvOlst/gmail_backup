package dropbox

import (
	"fmt"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

const name = "dropbox"

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

	// for _, entry := range res.Entries {
	// 	argm := files.NewGetMetadataArg("")

	// 	res, err := p.client.GetMetadata(argm)
	// 	// fmt.Printf("%+v\n", entry.IsMetadata)
	// 	// fmt.Println(entry.IsMetadata())
	// }

	// for res.HasMore {
	// 	arg := files.NewListFolderContinueArg(res.Cursor)

	// 	res, err = dbx.ListFolderContinue(arg)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	entries = append(entries, res.Entries...)
	// }

	// fmt.Println(r)
}

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {
	config := dropbox.Config{
		Token: cfg.AccessToken,
	}

	c := files.New(config)

	return &Provider{c}
}
