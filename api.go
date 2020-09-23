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
