package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	core "github.com/unf6/nucleus/internal/plugins"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:  "uninstall <pluginId>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dst := filepath.Join(core.InstallDir, args[0])
		if _, err := os.Stat(dst); err != nil {
			fmt.Println("Plugin not installed")
			return nil
		}

		if err := os.RemoveAll(dst); err != nil {
			return err
		}

		fmt.Println("Uninstalled plugin:", args[0])
		return nil
	},
}

func init() {
	Cmd.AddCommand(uninstallCmd)
}
