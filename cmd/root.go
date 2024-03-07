package cmd

import (
	"github.com/lancelee2885/personal-website-be/config"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	cobra.OnInitialize(config.InitializeConfig)

	root := &cobra.Command{
		Use: "pf-task",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	root.AddCommand(
		service(),
	)

	return root
}
