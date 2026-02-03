package shell

import (
	"os/exec"

	"github.com/unf6/nucleus/cmd/prompt"
)

// KillQuickShell forcibly stops any running QuickShell instances
func KillQuickShell() error {
	cmd := exec.Command("pkill", "-9", "quickshell")
	if err := cmd.Run(); err != nil {
		// Exit code 1 means no process was found → not an actual error
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil
		}
		prompt.Warn("Failed to kill Shell: " + err.Error())
		return err
	}
	prompt.Success("Shell stopped successfully")
	return nil
}
