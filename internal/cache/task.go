package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
	hash "github.com/segmentio/fasthash/fnv1a"
)

func ProjCachePath() (string, error) {
	global, err := util.GetGlobalConfigPath()
	if err != nil {
		return "", errors.Join(
			errors.New("Could not get global config path"),
			err,
		)
	}

	proj := hash.HashString64(util.ConfigPath)
	cache := filepath.Join(global, "cache", "task", hashStr(proj))

	// make sure the path exists
	if err := os.MkdirAll(cache, 0777); err != nil {
		return "", err
	}

	log.Debug("ProjCachePath", "cache", cache)

	return cache, nil
}

func LookupToolchain(in []string, out []string, tool string) (bool, error, func() error) {
	// cach file
	projCachePath, err := ProjCachePath()
	if err != nil {
		return false, errors.Join(
			errors.New("Failed to get project cache path"),
			err,
		), nil
	}
	if projCachePath == "" {
		return false, errors.New("Project cache path is empty"), nil
	}

	// check if tool has been run with these inputs and outputs
	toolchainHash := hash.Init64
	contentHash := hash.Init64
	for _, i := range in {
		toolchainHash = hash.AddString64(toolchainHash, i)
		if content, err := os.ReadFile(i); err != nil {
			return false, errors.Join(
				fmt.Errorf("Failed to read input file '%s'", i),
				err,
			), nil
		} else {
			contentHash = hash.AddBytes64(contentHash, content)
		}
	}
	for _, o := range out {
		toolchainHash = hash.AddString64(toolchainHash, o)
	}
	toolchainHash = hash.AddString64(toolchainHash, tool)
	contentHash = hash.AddString64(contentHash, tool)
	toolchainHashStr := hashStr(toolchainHash)
	toolchainCacheFile := filepath.Join(projCachePath, toolchainHashStr)

	fn := func() error {
		// write the new cache file
		if err := os.WriteFile(toolchainCacheFile, []byte(hashStr(contentHash)), 0777); err != nil {
			// _ = os.Remove(toolchainCacheFile)
			return errors.Join(
				errors.New("Failed to write toolchain cache file"),
				err,
			)
		}
		return nil
	}

	// check if every output file exists
	for _, o := range out {
		_, err := os.Stat(o)
		if err != nil {
			if os.IsNotExist(err) {
				log.Debug("Output file does not exist", "path", o)
				return false, nil, fn
			}
			return false, errors.Join(
				errors.New("Failed to check if output file exists"),
				err,
			), nil
		}
	}

	// check if every input file exists
	for _, i := range in {
		_, err := os.Stat(i)
		if err != nil {
			if os.IsNotExist(err) {
				return false, fmt.Errorf("Input file '%s' does not exist", i), nil
			}
			return false, errors.Join(
				errors.New("Failed to check if input file exists"),
				err,
			), nil
		}
	}

	stat, err := os.Stat(toolchainCacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("Toolchain cache file does not exist", "path", toolchainCacheFile)
			return false, nil, fn
		}
		return false, errors.Join(
			errors.New("Failed to check if toolchain cache file exists"),
			err,
		), nil
	}
	if stat.IsDir() {
		return false, fmt.Errorf("Toolchain cache file is a directory"), nil
	}

	// check if the toolchain cache file is up to date
	currentCacheContent, err := os.ReadFile(toolchainCacheFile)
	if err != nil {
		return false, errors.Join(
			errors.New("Failed to read toolchain cache file"),
			err,
		), nil
	}
	if hashStr(contentHash) == string(currentCacheContent) {
		log.Debug("Toolchain cache file is up to date", "path", toolchainCacheFile)
		return true, nil, fn
	} else {
		return false, nil, fn
	}
}

func hashStr(h uint64) string {
	return fmt.Sprintf("%x", h)
}
