package cmd

import (
        "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
        Use:   "nucleus",
        Short: "Nucleus Shell - Modern Wayland Shell Built For Hyprland",
        Long: `Nucleus Is A Beautiful, Customizable Shell.

Features:
  - Matugen Color Integration
  - Hot-reloading themes
  - Color Schemes
  - Fully customizable`,
}

func Execute() error {
        return rootCmd.Execute()
}
