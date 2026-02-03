package plugins

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/unf6/nucleus/cmd/prompt"
	"github.com/unf6/nucleus/internal/config"
)

var (
	InstallDir = filepath.Join(os.Getenv("HOME"), ".config/nucleus-shell/plugins")
	CacheBase  = "/tmp/nucleus-plugins"
)

// git runs a git command and returns an error if it fails
func git(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		prompt.Warn("Git command failed: git " + strings.Join(args, " "))
		return err
	}
	return nil
}

// UpdateAllRepos clones or updates all plugin repositories in cache
func UpdateAllRepos() error {
	for name, url := range config.Repos {
		dir := filepath.Join(CacheBase, name)

		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			prompt.Stage("Updating repository: " + name)
			if err := git("-C", dir, "pull", "--quiet"); err != nil {
				prompt.Warn("Failed to update repository: " + name)
				return err
			}
			continue
		}

		prompt.Stage("Cloning repository: " + name)
		_ = os.RemoveAll(dir)
		if err := git("clone", "--quiet", url, dir); err != nil {
			prompt.Warn("Failed to clone repository: " + name)
			return err
		}
	}
	prompt.Success("All plugin repositories updated successfully")
	return nil
}

// FindPlugin searches for a plugin by ID in cached repositories
func FindPlugin(id string) (repo string, path string, err error) {
	for repoName := range config.Repos {
		p := filepath.Join(CacheBase, repoName, id)
		if ValidatePluginDir(p) {
			return repoName, p, nil
		}
	}
	return "", "", errors.New("plugin not found: " + id)
}

// CopyDir copies a directory from src to dst
func CopyDir(src, dst string) error {
	if err := exec.Command("cp", "-r", src, dst).Run(); err != nil {
		prompt.Warn("Failed to copy directory: " + src + " → " + dst)
		return err
	}
	return nil
}
