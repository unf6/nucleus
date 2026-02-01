package plugins

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	core "github.com/unf6/nucleus/internal/plugins"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:  "install <pluginId>",
	Short: "Install A plugin",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := core.UpdateAllRepos(); err != nil {
			return err
		}

		repo, src, err := core.FindPlugin(args[0])
		_ = repo
		if err != nil {
			return err
		}

		m, err := core.LoadManifest(src + "/manifest.json")
		if err != nil {
			return err
		}

		if m.ID != args[0] {
			return errors.New("manifest id does not match plugin folder name")
		}

		dst := filepath.Join(core.InstallDir, args[0])
		if _, err := os.Stat(dst); err == nil {
			fmt.Println("Plugin already installed")
			return nil
		}

		_ = os.MkdirAll(core.InstallDir, 0755)
		if err := core.CopyDir(src, dst); err != nil {
			return err
		}

		fmt.Println("Installed plugin:", args[0])
		return nil
	},
}

func init() {
	Cmd.AddCommand(installCmd)
}
