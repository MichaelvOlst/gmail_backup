package dropbox

import (
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

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

	return false
}
