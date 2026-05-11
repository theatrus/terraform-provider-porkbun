package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type fileCredentials struct {
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

// loadCredentialsFile reads ~/.porkbun if present. A missing file is not
// an error; the returned struct is zero-valued in that case. A malformed
// file is surfaced as an error so misconfiguration is not silently ignored.
func loadCredentialsFile() (fileCredentials, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return fileCredentials{}, err
	}
	path := filepath.Join(home, ".porkbun")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fileCredentials{}, nil
		}
		return fileCredentials{}, fmt.Errorf("read %s: %w", path, err)
	}
	var creds fileCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return fileCredentials{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return creds, nil
}
