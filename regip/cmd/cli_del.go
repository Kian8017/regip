package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"regip"
)

var cliDelCmd = &cobra.Command{
	Use:   "del [id]",
	Short: "deletes a resource",
	Long:  `FIXME`,
	Args:  cobra.MinimumNArgs(1), // Only one resource type at a time! XD
	Run: func(cmd *cobra.Command, args []string) {
		lgr := CreateLogger("list", regip.CLR_cli)
		c, ok := CreateClient(lgr)
		if !ok {
			return
		}

		for _, a := range args {
			i, err := regip.ParseHex(a)
			if err != nil {
				lgr.Error("Couldn't parse ID: ", err)
				return
			}
			// Verify connectivity
			err = c.Delete(i)
			if err != nil {
				lgr.Error("Error deleting resource: ", err)
			} else {
				fmt.Println("Success")
			}
		}

		c.Wait()
	},
}

func init() {
	cliCmd.AddCommand(cliDelCmd)
}
