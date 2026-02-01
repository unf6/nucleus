package plugins

import (
	"encoding/json"
	"os"
)

type Manifest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Requires_Nucleus string `json:"requires_nucleus"`
}

func LoadManifest(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var m Manifest
	return &m, json.NewDecoder(f).Decode(&m)
}
