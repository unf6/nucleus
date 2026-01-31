package plugins

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage Nucleus Shell plugins",
}

func Init(root *cobra.Command) {
	root.AddCommand(Cmd)
}
