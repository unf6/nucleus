package cmd

import (
	"github.com/spf13/cobra"

	"github.com/unf6/nucleus/internal/ipc"
)

func ipcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ipc [target] [function] [args...]",
		Short: "Interact with the shell via IPC",
		Long: "Call QuickShell IPC targets and functions dynamically.",
		Args: cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			ipc.RunShellIPCCommand(args)
		},
	}
  
	cmd.ValidArgsFunction = func(
		_ *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		completions := ipc.GetShellIPCCompletions(args, toComplete)
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(ipcCmd())
}
