package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"regip"
)

var cliGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "gets a record",
	Long:  `FIXME`,
	Run: func(cmd *cobra.Command, args []string) {
		lgr := CreateLogger("list", regip.CLR_cli)
		c, ok := CreateClient(lgr)
		if !ok {
			lgr.Error("Couldn't create client")
			return
		}

		for _, a := range args {
			i, err := regip.ParseHex(a)
			if err != nil {
				lgr.Error("Couldn't parse ID: ", err)
				return
			}
			// Verify connectivity
			res, err := c.Get(i)
			if err != nil {
				lgr.Error("Error getting resource: ", err)
			} else {
				fmt.Println(res.String())
			}
		}

		c.Wait()
	},
}

func init() {
	cliCmd.AddCommand(cliGetCmd)
}
