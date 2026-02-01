package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	core "github.com/unf6/nucleus/internal/plugins"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:  "update <pluginId>",
  Short: "Update A Plugin",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := core.UpdateAllRepos(); err != nil {
			return err
		}

		dst := filepath.Join(core.InstallDir, args[0])
		if _, err := os.Stat(dst); err != nil {
			return fmt.Errorf("plugin not installed")
		}

		_, src, err := core.FindPlugin(args[0])
		if err != nil {
			return err
		}

		localM, _ := core.LoadManifest(dst + "/manifest.json")
		remoteM, _ := core.LoadManifest(src + "/manifest.json")

		if localM.Version == remoteM.Version {
			fmt.Println("Already up to date:", localM.Version)
			return nil
		}

		_ = os.RemoveAll(dst)
		if err := core.CopyDir(src, dst); err != nil {
			return err
		}

		fmt.Printf("Updated %s %s → %s\n",
			args[0], localM.Version, remoteM.Version)
		return nil
	},
}

func init() {
	Cmd.AddCommand(updateCmd)
}
