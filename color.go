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
	"github.com/gookit/color"
)

var (
	// Areas - Server
	CLR_server    = color.Cyan
	CLR_session   = color.Blue
	CLR_readloop  = color.Yellow
	CLR_writeloop = color.Red
	CLR_chainloop = color.White

	CLR_api = color.Green

	// Areas - Client
	CLR_cli = color.White

	// Types
	CLR_mt   = color.Magenta
	CLR_time = color.Yellow

	// Areas - Database
	CLR_db           = color.White
	CLR_fullindex    = color.Yellow
	CLR_indexrecords = color.Blue
)
