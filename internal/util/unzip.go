package util

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
)

// Unzip unzips an archive at its location
func Unzip(p string) error {
	log.Debug("Unzipping", "filepath", p)
	ext := filepath.Ext(p)
	switch ext {
	case ".zip":
		return unzip(p, filepath.Dir(p))
	case ".tar":
		r, err := os.Open(p)
		if err != nil {
			return err
		}
		defer r.Close()
		return untar(r, filepath.Dir(p))
	case ".gz":
		r, err := os.Open(p)
		if err != nil {
			return err
		}
		defer r.Close()
		gzipReader, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		return untar(gzipReader, filepath.Dir(p))
	case ".xz":
		return unxz(p, filepath.Dir(p))
	default:
		return fmt.Errorf("Unknown archive type '%s'", ext)
	}
}

// unzip unzips a zip archive at its location
func unzip(archPath string, outPath string) error {
	log.Debug("Unzipping", "outPath", outPath)
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error opening zip file '%s'", archPath),
			err,
		)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return errors.Join(
				fmt.Errorf("Error opening file '%s' in zip file '%s'", f.Name, archPath),
				err,
			)
		}
		defer rc.Close()

		p := filepath.Join(outPath, f.Name)

		if !IsInDir(p, outPath) {
			return fmt.Errorf("File '%s' is attempting to write outside of target directory", f.Name)
		}

		if err := os.MkdirAll(filepath.Dir(p), 0777); err != nil {
			return errors.Join(
				fmt.Errorf("Error creating directory '%s'", filepath.Dir(p)),
				err,
			)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		if f.Mode()&os.ModeSymlink != 0 {
			log.Error("Symlinks are not supported", "path", p, "archive", archPath)
			continue
		}

		log.Debug("Creating file", "path", p)

		file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()|0666)
		if err != nil {
			return errors.Join(
				fmt.Errorf("Error creating file '%s'", p),
				err,
			)
		}
		defer file.Close()

		if _, err = io.Copy(file, rc); err != nil {
			return errors.Join(
				fmt.Errorf("Error writing to file '%s'", p),
				err,
			)
		}
	}

	return nil
}

// untar untars a tar archive at its location
func untar(tarStream io.Reader, outPath string) error {
	log.Debug("Untarring", "outPath", outPath)
	tarReader := tar.NewReader(tarStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Join(
				fmt.Errorf("Error reading tar file"),
				err,
			)
		}

		p := filepath.Join(outPath, header.Name)

		if !IsInDir(p, outPath) {
			return fmt.Errorf("File '%s' is attempting to write outside of target directory", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(p, 0777); err != nil {
				return errors.Join(
					fmt.Errorf("Error creating directory '%s'", header.Name),
					err,
				)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(p), 0777); err != nil {
				return errors.Join(
					fmt.Errorf("Error creating directory '%s'", filepath.Dir(p)),
					err,
				)
			}
			outFile, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, header.FileInfo().Mode()|0666)
			if err != nil {
				return errors.Join(
					fmt.Errorf("Error creating file '%s'", p),
					err,
				)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.Join(
					fmt.Errorf("Error writing to file '%s'", p),
					err,
				)
			}
			// set executable bit
			if err := os.Chmod(p, header.FileInfo().Mode()); err != nil {
				return errors.Join(
					fmt.Errorf("Error setting executable bit on '%s'", p),
					err,
				)
			}
			outFile.Close()
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(p), 0777); err != nil {
				return errors.Join(
					fmt.Errorf("Error creating directory '%s'", filepath.Dir(p)),
					err,
				)
			}
			if err := os.Symlink(header.Linkname, p); err != nil {
				return errors.Join(
					fmt.Errorf("Error creating symlink '%s'", p),
					err,
				)
			}
		case tar.TypeLink:
			if err := os.MkdirAll(filepath.Dir(p), 0777); err != nil {
				return errors.Join(
					fmt.Errorf("Error creating directory '%s'", filepath.Dir(p)),
					err,
				)
			}
			if err := os.Link(header.Linkname, p); err != nil {
				return errors.Join(
					fmt.Errorf("Error creating hardlink '%s'", p),
					err,
				)
			}
		default:
			return errors.Join(
				fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, header.Name),
				err,
			)
		}
	}

	return nil
}

// unxz unxzs an .tar.xz archive at its location
func unxz(p string, outPath string) error {
	// is GNU tar available?
	if tar, err := exec.LookPath("tar"); err == nil {
		cmd := exec.Command(tar, "-xJf", p, "-C", outPath)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("GNU tar not found")
}
