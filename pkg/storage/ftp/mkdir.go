package ftp

import "fmt"

// Mkdir returns google_drive
func (p *Provider) Mkdir(path string) error {
	fmt.Println("TODO " + path)
	return nil
}

// IsNotExists check if a folder already exists
func (p *Provider) IsNotExists(err error) bool {
	return false
	// cerr, ok := err.(files.CreateFolderAPIError)
	// if !ok {
	// 	return false
	// }

	// if cerr.APIError.Error() == "path/conflict/folder/" {
	// 	return true
	// }
}
