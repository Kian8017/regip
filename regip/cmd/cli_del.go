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
