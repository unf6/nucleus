package cmd

import (
        "github.com/spf13/cobra"
        "github.com/unf6/nucleus/cmd/plugins"
)

var rootCmd = &cobra.Command{
        Use:   "nucleus",
        Short: "Nucleus Shell - A Modern Wayland Shell Built For Hyprland",
        Long: `Nucleus is a beautiful, customizable shell built to get things done`,
}

func Execute() error {
        return rootCmd.Execute()
}

func init() {
        plugins.Init(rootCmd)
}
