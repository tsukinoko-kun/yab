package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/mholt/archiver/v4"
)

var current = []string{"."}

// Unzip unzips an archive at its location
func Unzip(p string) error {
	target := filepath.Dir(p)
	log.Debug("Unzipping", "filepath", p, "target", target)
	fstream, err := os.Open(p)
	format, reader, err := archiver.Identify(p, fstream)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error opening archive '%s'", p),
			err,
		)
	}

	ctx := context.Background()

	if ex, ok := format.(archiver.Extractor); ok {
		ex.Extract(ctx, reader, current, func(ctx context.Context, f archiver.File) error {
			log.Debug("Extracting", "file", f.Name())
			return nil
		})
	}

	if decom, ok := format.(archiver.Decompressor); ok {
		rc, err := decom.OpenReader(reader)
		if err != nil {
			return errors.Join(fmt.Errorf("Error opening decompressor for '%s'", p), err)
		}
		defer rc.Close()

	}

	return nil
}
