package plugins

import (
	"github.com/charmbracelet/log"
  "github.com/unf6/nucleus/internal/config"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	InstallDir = filepath.Join(os.Getenv("HOME"), ".config/nucleus-shell/plugins")
	CacheBase  = "/tmp/nucleus-plugins"
)

func git(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func UpdateAllRepos() error {
	for name, url := range config.Repos {
		dir := filepath.Join(CacheBase, name)

		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			if err := git("-C", dir, "pull", "--quiet"); err != nil {
				return err
			}
			continue
		}

		_ = os.RemoveAll(dir)
		if err := git("clone", "--quiet", url, dir); err != nil {
			return err
		}
	}
	return nil
}

func FindPlugin(id string) (repo string, path string, err error) {
	for repo := range config.Repos {
		p := filepath.Join(CacheBase, repo, id)
		if ValidatePluginDir(p) {
			return repo, p, nil
		}
	}
	return "", "", log.Error("plugin not found")
}
