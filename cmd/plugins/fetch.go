package plugins

import (
	"fmt"
	"os"
	"path/filepath"

  "github.com/unf6/nucleus/internal/config"
	core "github.com/unf6/nucleus/internal/plugins"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch plugin metadata",
}

var fetchAllCmd = &cobra.Command{
	Use: "all",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := core.UpdateAllRepos(); err != nil {
			return err
		}

		for repo := range config.Repos {
			base := filepath.Join(core.CacheBase, repo)
			entries, _ := os.ReadDir(base)

			for _, e := range entries {
				dir := filepath.Join(base, e.Name())
				if !core.ValidatePluginDir(dir) {
					continue
				}

				m, err := core.LoadManifest(dir + "/manifest.json")
				if err != nil {
					continue
				}

				requires_nucleus := m.Requires_Nucleus
				if requires_nucleus == "" {
					requires_nucleus = "none"
				}

				fmt.Printf(
					"id: %s\nname: %s\nversion: %s\nauthor: %s\ndescription: %s\nrequires_nucleus: %s\nrepo: %s\n---\n",
					m.ID, m.Name, m.Version, m.Author,
					m.Description, m.Requires_Nucleus, repo,
				)
			}
		}
		return nil
	},
}

var fetchMachineCmd = &cobra.Command{
	Use: "all-machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := core.UpdateAllRepos(); err != nil {
			return err
		}

		for repo := range config.Repos {
			base := filepath.Join(core.CacheBase, repo)
			entries, _ := os.ReadDir(base)

			for _, e := range entries {
				dir := filepath.Join(base, e.Name())
				if !core.ValidatePluginDir(dir) {
					continue
				}

				m, _ := core.LoadManifest(dir + "/manifest.json")
				requires_nucleus := m.Requires_Nucleus
				if requires_nucleus == "" {
					requires_nucleus = "none"
				}

				fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					m.ID, m.Name, m.Version, m.Author,
					m.Description, m.Requires_Nucleus, repo,
				)
			}
		}
		return nil
	},
}

func init() {
	fetchCmd.AddCommand(fetchAllCmd, fetchMachineCmd)
	Cmd.AddCommand(fetchCmd)
}
