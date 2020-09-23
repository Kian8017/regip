/* Copyright (c) 2020 Kian Musser.
 * This file is part of regip.
 *
 * regip is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * regip is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with regip.  If not, see <https://www.gnu.org/licenses/>.
 */

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
