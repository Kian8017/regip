package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"regip"
)

var cliPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "verifies connectivity",
	Long:  `checks that we can connect to the server ok`,
	Run: func(cmd *cobra.Command, args []string) {
		lgr := CreateLogger("import", regip.CLR_cli)
		c, ok := CreateClient(lgr)
		if !ok {
			return
		}

		// Verify connectivity
		ok = c.Ping()
		if !ok {
			fmt.Println("connection failed")
		} else {
			fmt.Println("connection succeeded")
		}

		c.Wait()
	},
}

func init() {
	cliCmd.AddCommand(cliPingCmd)
}
