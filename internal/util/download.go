package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

func Download(url string, dest string) error {
	log.Debug("Downloading", "url", url, "dest", dest)
	resp, err := http.Get(url)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error downloading '%s'", url),
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		if location != "" {
			return Download(resp.Header.Get("Location"), dest)
		}
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error downloading '%s': %s", url, resp.Status)
	}

	// write to dest
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error creating file '%s'", dest),
			err,
		)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error writing to file '%s'", dest),
			err,
		)
	}

	return nil
}
