package downloader

import (
	"io"
	"net/http"
	"os"

	"github.com/LinkinStars/go-scaffold/logger"
)

func LoadAndSaveFile(url string, filePath string) error {
	logger.Debugf("trying to download file %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}
	return nil
}
