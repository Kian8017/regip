package cmd

import (
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "interact with a running regip server",
	Long:  `provides commands to interact with a running regip server`,
}

func init() {
	rootCmd.AddCommand(cliCmd)
}
