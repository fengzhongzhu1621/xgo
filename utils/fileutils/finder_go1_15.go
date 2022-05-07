//go:build !go1.16 || !finder
// +build !go1.16 !finder

package fileutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// FindConfigFile Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func FindConfigFile(fs afero.Fs, configPaths []string, configName string,
	supportedExts []string, configType string) (string, error) {
	log.Println("searching for config in paths", "paths", configPaths)

	for _, configPath := range configPaths {
		file := searchInPath(fs, configName, configPath, supportedExts, configType)
		if file != "" {
			return file, nil
		}
	}
	return "", ConfigFileNotFoundError{configName, fmt.Sprintf("%s", configPaths)}
}

func searchInPath(fs afero.Fs, configPath string, configName string,
	supportedExts []string, configType string) (filename string) {
	log.Println("searching for config in path", "path", configPath)
	for _, ext := range supportedExts {
		log.Println("checking if file exists", "file", filepath.Join(configPath, configName+"."+ext))
		if b, _ := Exists(fs, filepath.Join(configPath, configName+"."+ext)); b {
			log.Println("found file", "file", filepath.Join(configPath, configName+"."+ext))
			return filepath.Join(configPath, configName+"."+ext)
		}
	}

	if configType != "" {
		if b, _ := Exists(fs, filepath.Join(configPath, configName)); b {
			return filepath.Join(configPath, configName)
		}
	}

	return ""
}

// Check if file Exists
func Exists(fs afero.Fs, path string) (bool, error) {
	stat, err := fs.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
