package cmd

import (
//	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/unf6/nucleus/internal/config"
	"github.com/unf6/nucleus/internal/shell"
   	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start The Nucleus Shell",
	Long: `Start the Nucleus Shell using QuickShell.

This will launch QuickShell with the Nucleus Shell configuration.`,
	RunE: runShell,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("reload", "r", false, "Kill existing instance before starting")
	runCmd.Flags().BoolP("debug", "d", false, "Run in foreground (without --daemonize)")
}

func runShell(cmd *cobra.Command, args []string) error {
	reload, _ := cmd.Flags().GetBool("reload")
	debug, _ := cmd.Flags().GetBool("debug")

	if !config.IsInstalled() {
		log.Error("Nucleus Shell is not installed. Run 'nucleus install' first")
	}

	if reload {
		log.Info("Stopping existing QuickShell instances...")
		if err := shell.KillQuickShell(); err != nil {
			log.Error("Warning: %v\n", err)
		}
	}

	shellPath, err := config.GetConfigDir()
	if err != nil {
		 log.Error("failed to get shell config path: %w", err)
	}

	var quickshellCmd *exec.Cmd
	if debug {
		log.Debug("Starting Nucleus Shell in debug mode (foreground)...")
		quickshellCmd = exec.Command("quickshell", "--no-duplicate", "-c", shellPath)
	} else {
		log.Info("Starting Nucleus Shell.")
		quickshellCmd = exec.Command("quickshell", "--no-duplicate", "--daemonize", "-c", shellPath)
	}
	quickshellCmd.Stdout = os.Stdout
	quickshellCmd.Stderr = os.Stderr

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := quickshellCmd.Start(); err != nil {
		 log.Error("failed to start QuickShell: %w", err)
	}

	log.Warn("Nucleus Shell started (PID: %d)\n", quickshellCmd.Process.Pid)

	go func() {
		<-sigChan
		log.Error("\nShutting down Nucleus Shell...")
		quickshellCmd.Process.Signal(syscall.SIGTERM)
	}()

	if err := quickshellCmd.Wait(); err != nil {
		 log.Error("QuickShell exited with error", err)
	}

	return nil
}
