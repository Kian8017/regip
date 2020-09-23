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
