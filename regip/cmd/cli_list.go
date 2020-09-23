package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"regip"
)

var cliListCmd = &cobra.Command{
	Use:   "list [resource type]",
	Short: "lists records",
	Long:  `FIXME`,
	Args:  cobra.ExactArgs(1), // Only one resource type at a time! XD
	Run: func(cmd *cobra.Command, args []string) {
		lgr := CreateLogger("list", regip.CLR_cli)
		rt, ok := regip.StringToResourceType(args[0])
		if !ok {
			lgr.Error("Invalid resource type '", args[0], "'")
			return
		}
		c, ok := CreateClient(lgr)
		if !ok {
			return
		}

		resCh := c.List(rt)
		total := 0
		for cur := range resCh {
			if onlyID {
				fmt.Println(cur.ID().String())
			} else {
				fmt.Println(cur.String())
			}
			total++
		}
		if !onlyID {
			fmt.Println("Total: ", total)
		}

		c.Wait()
	},
}

func init() {
	cliCmd.AddCommand(cliListCmd)
}
