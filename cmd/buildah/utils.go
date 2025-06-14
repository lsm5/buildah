package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// createTempStorageConf creates a temporary storage.conf file with the specified digest type
func createTempStorageConf(digestType string) (string, error) {
	// Find the current storage.conf
	storageConfPath := os.Getenv("CONTAINERS_STORAGE_CONF")
	if storageConfPath == "" {
		// Try common locations
		candidates := []string{
			"/etc/containers/storage.conf",
			"/usr/share/containers/storage.conf",
			filepath.Join(os.Getenv("HOME"), ".config/containers/storage.conf"),
		}
		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				storageConfPath = candidate
				break
			}
		}
	}
	if storageConfPath == "" {
		return "", fmt.Errorf("could not find storage.conf to override digest_type")
	}

	// Read and modify storage.conf
	orig, err := ioutil.ReadFile(storageConfPath)
	if err != nil {
		return "", fmt.Errorf("failed to read storage.conf: %w", err)
	}
	conf := string(orig)

	// Replace or add digest_type in [storage.options]
	re := regexp.MustCompile(`(?m)^([ \t]*#?[ \t]*digest_type[ \t]*=[ \t]*).*$`)
	if re.MatchString(conf) {
		conf = re.ReplaceAllString(conf, "digest_type = \""+digestType+"\"")
	} else {
		// Add under [storage.options]
		reSection := regexp.MustCompile(`(?m)^\[storage.options\]$`)
		if reSection.MatchString(conf) {
			conf = reSection.ReplaceAllString(conf, "[storage.options]\ndigest_type = \""+digestType+"\"")
		} else {
			// Add section if missing
			conf += "\n[storage.options]\ndigest_type = \"" + digestType + "\"\n"
		}
	}

	// Write to temp file
	tmpFile, err := ioutil.TempFile("", "storage-*.conf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp storage.conf: %w", err)
	}
	tempStorageConf := tmpFile.Name()
	if _, err := tmpFile.Write([]byte(conf)); err != nil {
		tmpFile.Close()
		os.Remove(tempStorageConf)
		return "", fmt.Errorf("failed to write temp storage.conf: %w", err)
	}
	tmpFile.Close()

	return tempStorageConf, nil
}
