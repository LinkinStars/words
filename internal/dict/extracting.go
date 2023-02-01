package dict

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/words/assets"
	"github.com/LinkinStars/words/internal/config"
)

func extractBooks() (err error) {
	logger.Debug("try to extract books...")
	booksData, err := assets.Books.ReadFile("books.tar.gz")
	if err != nil {
		return err
	}
	gr, err := gzip.NewReader(bytes.NewBuffer(booksData))
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		if strings.HasPrefix(hdr.Name, ".") {
			continue
		}
		filename := filepath.Join(config.BookDir, hdr.Name)
		logger.Debugf("try to extract book %s", filename)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
	}
	return nil
}
