package plugins

import "os"

func ValidatePluginDir(dir string) bool {
	required := []string{
		"Main.qml",
		"PluginConfigData.qml",
		"Settings.qml",
		"manifest.json",
	}

	for _, f := range required {
		if _, err := os.Stat(dir + "/" + f); err != nil {
			return false
		}
	}
	return true
}
