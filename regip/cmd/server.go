package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"regip"
	"strings"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server [path to database]",
	Args:  cobra.MinimumNArgs(1), // FIXME: also check for path existence / database existence
	Short: "Start the server",
	Long: `server starts the built-in server.
Add usage examples here f.ex (:2020)`,
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: Don't hardcode
		addr := regip.DEFAULT_LISTEN_ADDR
		databaseLoc := strings.Join(args, " ")
		// FIXME: add -v flag, rather than logging by default
		serv, err := regip.NewServer(databaseLoc, addr, true)
		if err != nil {
			// FIXME: add error formatting
			fmt.Println(err)
			return
		}
		defer serv.Close()

		serv.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
