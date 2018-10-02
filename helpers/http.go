package helpers

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	ContentTypeApplicationJson = "application/json"
)

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Printf("[DEBUG] Donwloading... FilePath: %s URL: %s\n", filepath, url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
