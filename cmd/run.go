package cmd

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/unf6/nucleus/internal/config"
	"github.com/unf6/nucleus/internal/shell"
	"github.com/unf6/nucleus/internal/prompt"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Nucleus Shell",
	Long: `Start the Nucleus Shell using QuickShell.

This will launch QuickShell with the Nucleus Shell configuration.`,
	RunE: runShell,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("reload", "r", false, "Stop existing instance before starting")
	runCmd.Flags().BoolP("debug", "d", false, "Run in foreground (no --daemonize)")
}

func runShell(cmd *cobra.Command, args []string) error {
	reload, _ := cmd.Flags().GetBool("reload")
	debug, _ := cmd.Flags().GetBool("debug")

	if !config.IsInstalled() {
		return prompt.Fail("Nucleus Shell is not installed. Run 'nucleus install' first.")
	}

	if reload {
		prompt.Stage("Stopping existing QuickShell instances...")
		if err := shell.KillQuickShell(); err != nil {
			prompt.Warn("Warning: %v", err)
		}
	}

	shellPath, err := config.GetConfigDir()
	if err != nil {
		return prompt.Fail("Failed to get shell config path: %v", err)
	}

	var quickshellCmd *exec.Cmd
	if debug {
		prompt.Stage("Starting Nucleus Shell in debug mode (foreground)...")
		quickshellCmd = exec.Command("quickshell", "--no-duplicate", "-c", shellPath)
	} else {
		prompt.Stage("Starting Nucleus Shell in background...")
		quickshellCmd = exec.Command("quickshell", "--no-duplicate", "--daemonize", "-c", shellPath)
	}
	quickshellCmd.Stdout = os.Stdout
	quickshellCmd.Stderr = os.Stderr

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := quickshellCmd.Start(); err != nil {
		return prompt.Fail("Failed to start QuickShell: %v", err)
	}

	prompt.Success("Nucleus Shell started (PID: %d)", quickshellCmd.Process.Pid)

	go func() {
		<-sigChan
		prompt.Warn("\nShutting down Nucleus Shell...")
		_ = quickshellCmd.Process.Signal(syscall.SIGTERM)
	}()

	if err := quickshellCmd.Wait(); err != nil {
		prompt.Warn("QuickShell exited with error: %v", err)
	}

	return nil
}
