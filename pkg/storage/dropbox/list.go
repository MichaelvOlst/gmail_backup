package dropbox

import (
	"fmt"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

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
