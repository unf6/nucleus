package shell

import (
        "github.com/charmbracelet/log"                                                       
        "os/exec"
)

func KillQuickShell() error {
        cmd := exec.Command("pkill", "-f", "quickshell")
        if err := cmd.Run(); err != nil {
                if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
                        return nil                             
                }
                return log.Error("Failed To Kill QuickShell: %w", err)
        }
        return nil
}
