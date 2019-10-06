package ftp

import (
	"fmt"
	"io"
	"os"
)

// Put returns google_drive
func (p *Provider) Put(filename, path string, file *os.File, r io.Reader) error {
	fmt.Println("TODO " + filename)
	return nil
}
