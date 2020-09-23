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
	"fmt"
)

type API struct {
	database  *DB
	endpoints map[MessageType]Endpoint
}

func NewAPI(d *DB) *API {
	var a API
	a.database = d
	a.endpoints = ENDPOINTS
	fmt.Println(fmt.Sprintf("DEBUG: API has %d callbacks registered", len(a.endpoints)))
	return &a
}

// The mastermind -- delegates to sub handlers
func (a *API) Handle(first *Message, input, output chan *Message, userid *string, lgr *Logger) chan struct{} {
	lgr = lgr.Tag("api", CLR_api)
	// Check if we have that type implemented,
	// Then check auth,
	// Then call
	// Or scream "Fire!" in a crowded theater because we don't have a callback implemented for that Call

	// Loggertag is 'call id'
	lgTag := fmt.Sprintf("%s %d", first.Type, first.Id)
	callLogger := lgr.Tag(lgTag, CLR_mt)
	callLogger.Print("called")

	cb, ok := a.endpoints[first.Type]
	if !ok {
		// Here's where you scream "Fire!"
		return Unknown(a.database, userid, first, input, output, callLogger)
	}
	// We made it!
	// FIXME Call
	return cb(a.database, userid, first, input, output, callLogger)
}

var Unknown = WrapEndpoint(func(_ *DB, _ *string, first *Message, inp, oup chan *Message, lgr *Logger) {
	lgr.Error("unimplemented")
	oup <- NewMessage(first.Id, MT_errnotimplemented, "")
})
