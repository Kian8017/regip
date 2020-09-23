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
