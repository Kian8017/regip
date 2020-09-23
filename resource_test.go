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

package regip_test

import (
	"regip"
	"testing"
)

func TestInterfaces(t *testing.T) {
	var _ regip.Resource = &regip.Record{}
	var _ regip.Resource = &regip.User{}
	var _ regip.Resource = regip.Country{}
	var _ regip.Resource = regip.IndexRecord{}
	var _ regip.Resource = &regip.Trigram{}
}
