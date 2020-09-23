package cmd

import (
	"github.com/spf13/cobra"
	"regip"
)

var cliFullIndexCmd = &cobra.Command{
	Use:   "fullindex",
	Short: "runs a full index",
	Long:  `FIXME`,
	Run: func(cmd *cobra.Command, args []string) {
		lgr := CreateLogger("fullindex", regip.CLR_cli)
		c, ok := CreateClient(lgr)
		if !ok {
			lgr.Error("Couldn't create client")
			return
		}

		err := c.FullIndex()
		if err != nil {
			lgr.Error(err)
		}

		c.Wait()
	},
}

func init() {
	cliCmd.AddCommand(cliFullIndexCmd)
}
