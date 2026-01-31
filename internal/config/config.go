package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigDirName = "nucleus-shell"
	RepoURL       = "https://github.com/xzepyx/nucleus-shell.git"
)

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "quickshell", ConfigDirName), nil
}

func IsInstalled() bool {
	configDir, err := GetConfigDir()
	if err != nil {
		return false
	}
	
	shellFile := filepath.Join(configDir, "shell.qml")
	_, err = os.Stat(shellFile)
	return err == nil
}

func GetShellFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "shell.qml"), nil
}
