package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	fmt.Println("Updating Nucleus Shell...")

	gitCmd := exec.Command("git", "-C", configDir, "pull", "origin", "main")
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return err
	}
	fmt.Println("\n✅ Updated successfully!")
	return nil
}
