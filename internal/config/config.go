package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigDirName = "nucleus-shell"
	RepoURL       = "https://github.com/xzepyx/nucleus-shell.git"
)

var Repos = map[string]string{
	"official": "https://github.com/xZepyx/nucleus-plugins.git",
}


func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "quickshell", ConfigDirName), nil
}

func GetPluginDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", ConfigDirName, "plugins"), nil
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
