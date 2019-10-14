package ftp

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Put stores the file on the ftp location
func (p *Provider) Put(filename, path string, file *os.File, r io.Reader) error {

	path = strings.TrimRight(path, "/")
	filenameBase := fmt.Sprintf("%s/%s", path, filename)

	err := p.client.Stor(filenameBase, r)
	if err != nil {
		return err
	}
	return nil
}
