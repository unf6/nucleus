package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"archive/zip"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/unf6/nucleus/internal/config"
	"github.com/unf6/nucleus/internal/installer"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or Update Nucleus Shell",
	RunE:  runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("force", "f", false, "Force reinstall (removes existing installation)")
	installCmd.Flags().Bool("dev", false, "Clone full git repository (for development)")
    installCmd.Flags().Bool("stable", false, "Update latest stable release")
    installCmd.Flags().Bool("indev", false, "Update latest pre-release")

}

func runInstall(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	if config.IsInstalled() && !force {
		fmt.Println("Nucleus Shell is already installed.")
		fmt.Print("Would you like to update instead? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(strings.ToLower(resp))

		if resp == "y" || resp == "yes" {
			return runUpdate()
		}

		fmt.Println("Use --force to reinstall.")
		return nil
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	if force && config.IsInstalled() {
		fmt.Println("Removing existing installation...")
		if err := os.RemoveAll(configDir); err != nil {
			return err
		}
	}
	
	fmt.Println("Installing Nucleus Shell...")
	fmt.Printf("Cloning into %s\n", configDir)

	gitCmd := exec.Command("git", "clone", config.RepoURL, configDir)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if err := gitCmd.Run(); err != nil {
		return err
	}
	
if err := installer.RunWithSpinner(
	"Installing The Dependencies",
	installer.InstallDependencies,
); err != nil {
	return err
}

if err := installer.RunWithSpinner(
	"Copying Nucleus Shell to QuickShell config",
	installer.CopyToQuickShellConfig,
); err != nil {
	return err
}

	fmt.Println("\n✓ Nucleus shell installed successfully!")
	fmt.Println("Run it with:")
	fmt.Println("  nucleus run")

	return nil
}

  func runUpdate() error {
	const (
		repo = "xZepyx/nucleus-shell"
		api  = "https://api.github.com/repos/" + repo + "/releases"
	)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".config", "nucleus-shell", "config", "configuration.json")
	qsDir := filepath.Join(home, ".config", "quickshell", "nucleus-shell")

	stable, _ := installCmd.Flags().GetBool("stable")
	indev, _ := installCmd.Flags().GetBool("indev")

	if stable && indev {
		return fmt.Errorf("cannot use --stable and --indev together")
	}

	mode := "stable"
	if indev {
		mode = "indev"
	}

	cfgRaw, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("configuration.json not found")
	}

	var cfg map[string]any
	if err := json.Unmarshal(cfgRaw, &cfg); err != nil {
		return err
	}

	shell, ok := cfg["shell"].(map[string]any)
	if !ok {
		return fmt.Errorf("shell.version not found")
	}

	current, ok := shell["version"].(string)
	if !ok || current == "" {
		return fmt.Errorf("shell.version not found")
	}

	resp, err := http.Get(api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var releases []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return err
	}

	var latestTag string
	for i := len(releases) - 1; i >= 0; i-- {
		r := releases[i]

		if r["draft"].(bool) {
			continue
		}

		if mode == "stable" && r["prerelease"].(bool) {
			continue
		}

		latestTag, _ = r["tag_name"].(string)
		break
	}

	if latestTag == "" {
		return fmt.Errorf("failed to determine release version")
	}

	latest := strings.TrimPrefix(latestTag, "v")

	if latest == current {
		fmt.Printf("already up to date (%s)\n", current)
		return nil
	}

	tmpDir, err := os.MkdirTemp("", "nucleus-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "source.zip")
	zipURL := fmt.Sprintf("https://github.com/%s/archive/refs/tags/%s.zip", repo, latestTag)

	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err = http.Get(zipURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(tmpDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		dst, err := os.Create(path)
		if err != nil {
			return err
		}

		src, err := f.Open()
		if err != nil {
			dst.Close()
			return err
		}

		_, err = io.Copy(dst, src)
		dst.Close()
		src.Close()

		if err != nil {
			return err
		}
	}

	srcDir := filepath.Join(
		tmpDir,
		fmt.Sprintf("nucleus-shell-%s", latest),
		"quickshell",
		"nucleus-shell",
	)

	if _, err := os.Stat(srcDir); err != nil {
		return fmt.Errorf("nucleus-shell folder not found in source archive")
	}

	os.RemoveAll(qsDir)
	if err := os.MkdirAll(qsDir, 0755); err != nil {
		return err
	}

	if err := exec.Command("cp", "-r", srcDir+string(os.PathSeparator)+".", qsDir).Run(); err != nil {
		return err
	}

	shell["version"] = latest
	cfg["shell"] = shell

	updated, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(configPath, updated, 0644); err != nil {
		return err
	}

	exec.Command("killall", "qs").Run()
	exec.Command("nohup", "qs", "-c", "nucleus-shell").Start()

	fmt.Printf("Updated %s → %s\n", current, latest)
	return nil
}
