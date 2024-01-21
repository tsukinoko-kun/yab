package util

import (
	"os"
	"path/filepath"
)

// Find finds a file in a directory
func Find(name string, in string) (string, error) {
	results := []string{}
	err := filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == name {
			results = append(results, path)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", os.ErrNotExist
	}

	// return the shortest path
	shortest := results[0]
	length := len(shortest)
	for _, result := range results {
		if len(result) < length {
			shortest = result
			length = len(result)
		}
	}

	return shortest, nil
}
