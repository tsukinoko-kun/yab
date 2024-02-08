package zip

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"

	lua "github.com/Frank-Mayer/yab/internal/lua"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

var (
	// array of filenames to exclude from zip
	zipBlacklist = [...]string{
		".DS_Store",
		"thumbs.db",
	}
)

// Create a zip file containing the given files. Returns true if successful, false otherwise.
func MakeZip(l *lua.LState) int {
	filesT := l.CheckTable(1)
	output := l.CheckString(2)

	cancelSpinner := util.Spin()
	defer cancelSpinner()

	archive, err := os.Create(output)
	if err != nil {
		l.Error(lua.LString("Error creating zip file. "+err.Error()), 0)
		return 0
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)
	defer writer.Close()

	filesT.ForEach(func(_ lua.LValue, value lua.LValue) {
		path := value.String()

		fi, err := os.Stat(path)
		if os.IsNotExist(err) {
			l.Error(lua.LString("File does not exist. "+err.Error()), 0)
			return // continue
		}

		if fi.IsDir() {
			if err := addFilesToZip(writer, path, path); err != nil {
				l.Error(lua.LString("Error adding files to zip. "+err.Error()), 0)
				return // continue
			}
		} else {
			dat, err := os.ReadFile(path)
			if err != nil {
				l.Error(lua.LString("Error reading file. "+err.Error()), 0)
				return // continue
			}

			f, err := writer.Create(fi.Name())
			if err != nil {
				l.Error(lua.LString("Error creating file in zip. "+err.Error()), 0)
				return // continue
			}
			_, err = f.Write(dat)
			if err != nil {
				l.Error(lua.LString("Error writing file to zip. "+err.Error()), 0)
				return // continue
			}
		}
	})

	return 0
}

func addFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

filesLoop:
	for _, file := range files {
		// check if file is in blacklist
		for _, blacklisted := range zipBlacklist {
			if file.Name() == blacklisted {
				continue filesLoop
			}
		}

		fullfilepath := filepath.Join(basePath, file.Name())
		fi, err := os.Stat(fullfilepath)
		if os.IsNotExist(err) {
			continue
		}

		if fi.Mode()&(fs.ModeIrregular|fs.ModeSymlink|fs.ModeDevice|fs.ModeNamedPipe|fs.ModeSocket) != 0 {
			// skip irregular files (e.g. symlinks)
			log.Warn("Skipping irregular file", "path", fullfilepath)
			continue
		} else if fi.IsDir() {
			if err := addFilesToZip(w, fullfilepath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else {
			inZipPath := filepath.Join(baseInZip, file.Name())
			dat, err := os.ReadFile(fullfilepath)
			if err != nil {
				return err
			}

			f, err := w.Create(inZipPath)
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
