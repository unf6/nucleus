package config

import (
	"os"
	"path/filepath"
)

const (
	// ConfigDirName is the directory name for Nucleus Shell config
	ConfigDirName = "nucleus-shell"

	// RepoURL is the default repository URL for Nucleus Shell
	RepoURL = "https://github.com/xzepyx/nucleus-shell.git"
)

// Repos defines plugin repositories
var Repos = map[string]string{
	"official": "https://github.com/xZepyx/nucleus-plugins.git",
}

// GetConfigDir returns the full path to the Nucleus Shell config directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "quickshell", ConfigDirName), nil
}

// IsInstalled checks if the shell.qml file exists in the config directory
func IsInstalled() bool {
	configDir, err := GetConfigDir()
	if err != nil {
		return false
	}
	shellFile := filepath.Join(configDir, "shell.qml")
	_, err = os.Stat(shellFile)
	return err == nil
}

// GetShellFile returns the full path to the shell.qml file
func GetShellFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "shell.qml"), nil
}
