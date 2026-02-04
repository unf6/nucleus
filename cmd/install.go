package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/unf6/nucleus/internal/config"
	"github.com/unf6/nucleus/internal/installer"
	"github.com/unf6/nucleus/cmd/prompt"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update Nucleus Shell",
	RunE:  runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("force", "f", false, "Force reinstall (removes existing installation)")
	installCmd.Flags().BoolP("source", "s", false, "Clone the full git repository instead of using latest tag")
	installCmd.Flags().String("git", "", "Install or update from a specific git tag")
}

func runInstall(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	source, _ := cmd.Flags().GetBool("source")
	gitTag, _ := cmd.Flags().GetString("git")

	mode := "latest-tag" // default
	if source {
		mode = "source"
	}
	if gitTag != "" {
		mode = "git:" + gitTag
	}

	// If installed and not forced, prompt update
	if config.IsInstalled() && !force {
		prompt.Stage("Nucleus Shell is already installed.")
		resp := prompt.Ask("Update instead? [y/N]: ")

		if resp == "y" || resp == "yes" {
			return runUpdate(mode)
		}
		prompt.Action("Use --force to reinstall.")
		return nil
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to get config directory: %v", err))
	}

	if force && config.IsInstalled() {
		prompt.Stage("Removing existing installation...")
		if err := os.RemoveAll(configDir); err != nil {
			return prompt.Fail(fmt.Sprintf("Failed to remove existing installation: %v", err))
		}
	}

	if mode == "source" {
		// Full repo clone
		prompt.Stage("Installing from full git repository...")
		prompt.Action(fmt.Sprintf("Cloning into %s", configDir))

		gitCmd := exec.Command("git", "clone", config.RepoURL, configDir)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			return prompt.Fail(fmt.Sprintf("Git clone failed: %v", err))
		}
	} else {
		// Install from latest tag or specified tag
		if err := installFromTag(mode, configDir); err != nil {
			return err
		}
	}

	if err := installer.RunWithSpinner("Installing Dependencies", installer.InstallDependencies); err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to install dependencies: %v", err))
	}

	if err := installer.RunWithSpinner("Copying Nucleus Shell to QuickShell config", installer.CopyToQuickShellConfig); err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to copy shell to QuickShell: %v", err))
	}

	prompt.Success("Nucleus Shell installed successfully!")
	prompt.Action("Run it with: nucleus run")
	return nil
}

// runUpdate updates existing installation (downloads latest tag or git tag)
func runUpdate(mode string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to get home directory: %v", err))
	}

	configPath := filepath.Join(home, ".config", "nucleus-shell", "config", "configuration.json")
	cfgRaw, err := os.ReadFile(configPath)
	if err != nil {
		return prompt.Fail("configuration.json not found")
	}

	var cfg map[string]any
	if err := json.Unmarshal(cfgRaw, &cfg); err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to parse configuration.json: %v", err))
	}

	shellCfg, ok := cfg["shell"].(map[string]any)
	if !ok {
		return prompt.Fail("shell.version not found in configuration")
	}

	current, _ := shellCfg["version"].(string)
	if current == "" {
		current = "unknown"
	}

	qsDir := filepath.Join(home, ".config", "quickshell", "nucleus-shell")

	prompt.Stage(fmt.Sprintf("Updating Nucleus Shell (current version: %s)", current))
	if err := installFromTag(mode, qsDir); err != nil {
		return err
	}

	// Update version in config
	if strings.HasPrefix(mode, "git:") {
		shellCfg["version"] = strings.TrimPrefix(mode, "git:")
	} else {
		shellCfg["version"] = "latest"
	}
	cfg["shell"] = shellCfg

	updated, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(configPath, updated, 0644)

	_ = exec.Command("killall", "quickshell").Run()
	_ = exec.Command("nohup", "quickshell", "-c", "nucleus-shell").Start()

	prompt.Success("Update complete!")
	return nil
}

// installFromTag downloads a release zip (latest or git tag) and installs it
func installFromTag(mode, targetDir string) error {
	const repo = "xZepyx/nucleus-shell"
	api := "https://api.github.com/repos/" + repo + "/releases"

	var tag string
	if strings.HasPrefix(mode, "git:") {
		tag = strings.TrimPrefix(mode, "git:")
		tag = strings.TrimPrefix(tag, "v")
		if tag == "" {
			return prompt.Fail("Invalid git tag specified")
		}
		tag = "v" + tag
	} else {
		prompt.Stage("Fetching latest release info from GitHub...")
		resp, err := http.Get(api)
		if err != nil {
			return prompt.Fail(fmt.Sprintf("Failed to query GitHub API: %v", err))
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return prompt.Fail(fmt.Sprintf("GitHub API returned status %d", resp.StatusCode))
		}

		var releases []struct {
			TagName    string `json:"tag_name"`
			Draft      bool   `json:"draft"`
			Prerelease bool   `json:"prerelease"`
			Published  string `json:"published_at"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			return prompt.Fail(fmt.Sprintf("Failed to decode releases: %v", err))
		}

		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Published < releases[j].Published
		})

		for i := len(releases) - 1; i >= 0; i-- {
			r := releases[i]
			if r.Draft || r.TagName == "" {
				continue
			}
			tag = r.TagName
			break
		}

		if tag == "" {
			return prompt.Fail("Failed to determine latest release version")
		}
	}

	prompt.Stage(fmt.Sprintf("Installing version %s", tag))

	tmpDir, err := os.MkdirTemp("", "nucleus-install-*")
	if err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to create temp directory: %v", err))
	}
	defer os.RemoveAll(tmpDir)

	zipURL := fmt.Sprintf("https://github.com/%s/archive/refs/tags/%s.zip", repo, tag)
	zipPath := filepath.Join(tmpDir, "source.zip")

	prompt.Stage("Downloading release archive...")
	resp, err := http.Get(zipURL)
	if err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to download archive: %v", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return prompt.Fail(fmt.Sprintf("Failed to download archive: status %d", resp.StatusCode))
	}

	out, err := os.Create(zipPath)
	if err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to create zip file: %v", err))
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return prompt.Fail(fmt.Sprintf("Failed to save zip file: %v", err))
	}
	out.Close()

	prompt.Stage("Extracting archive...")
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to open zip: %v", err))
	}
	defer zr.Close()

	for _, f := range zr.File {
		path := filepath.Join(tmpDir, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, 0755)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(path), 0755)
		dst, _ := os.Create(path)
		src, _ := f.Open()
		_, _ = io.Copy(dst, src)
		dst.Close()
		src.Close()
	}

	srcDir := filepath.Join(tmpDir, fmt.Sprintf("nucleus-shell-%s", strings.TrimPrefix(tag, "v")), "quickshell", "nucleus-shell")
	if _, err := os.Stat(srcDir); err != nil {
		return prompt.Fail("nucleus-shell directory missing in archive")
	}

	prompt.Stage("Installing files...")
	_ = os.RemoveAll(targetDir)
	_ = os.MkdirAll(targetDir, 0755)
	if err := exec.Command("cp", "-r", srcDir+string(os.PathSeparator)+".", targetDir).Run(); err != nil {
		return prompt.Fail(fmt.Sprintf("Failed to copy files: %v", err))
	}

	return nil
}
