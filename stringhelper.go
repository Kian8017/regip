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

package regip

import (
	"strings"
)

const DROP string = "\r\n\ufeff" // Drop Byte Order Mark
const REPLACE string = "\t"

func NormalizeString(s string) string {
	ns := strings.ToLower(s)
	var sb strings.Builder
	for _, r := range ns {
		if IsInList(r, REPLACE) {
			sb.WriteString(" ")
		} else if IsInList(r, DROP) {
			continue
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func GenerateTrigrams(s string, pad bool) []string {
	var ns string
	if pad {
		ns = NormalizeString("  " + s + "  ")
	} else {
		ns = NormalizeString(s)
	}
	var ret []string
	for i := 0; i < len(ns)-2; i++ {
		ret = append(ret, ns[i:i+3])
	}
	return ret
}

func IsInList(c rune, list string) bool {
	for _, r := range list {
		if r == c {
			return true
		}
	}
	return false
}
