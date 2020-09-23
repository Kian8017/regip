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

func FilterName(want string, inputF *Flow) *Flow {
	return NewFlowFilter(func(n Resource) bool {
		res, ok := n.(Nameable)
		if !ok {
			return false
		}
		return strings.Contains(res.Name(), want)
	})(inputF)
}

/*
func FilterType(nt RecordType, inputF *Flow) *Flow {
	return NewFlowFilter(func(n *Record) bool {
		return n.Type == nt
	})(inputF)
}

func FilterCountry(ci ID, inputF *Flow) *Flow {
	return NewFlowFilter(func(n *Record) bool {
		return n.Country.Equal(ci)
	})(inputF)
}
*/
