package gmail

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"google.golang.org/api/gmail/v1"
)

func (g *Gmail) saveMessage(path string, m *gmail.Message, zipWriter *zip.Writer) error {

	b, _ := g.decodeMessage(m)
	t := time.Unix(m.InternalDate/1000, 0)
	folder := fmt.Sprintf("%s", t.Format("2006-01"))

	tempFile, err := ioutil.TempFile("", "backup_")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(b)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s.eml", folder, m.Id)

	info, err := tempFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	writer.Write(b)
	_, err = io.Copy(writer, tempFile)
	return err
}

func (g *Gmail) decodeMessage(m *gmail.Message) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(m.Raw)
}
