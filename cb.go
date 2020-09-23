package regip

import (
	"encoding/json"
	"time"
)

var ENDPOINTS map[MessageType]Endpoint

// SECTION:AUTH

type Userpass struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

var HandleLogin = WrapEndpoint(func(d *DB, userid *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	var up Userpass
	err := json.Unmarshal([]byte(req.Payload), &up)
	if err != nil {
		lgr.Error("invalid formatting", err)
		output <- NewMessage(req.Id, MT_errinvalidformatting, "")
		return
	}

	lgr.Printf("%s is trying to login", up.User)
	// Get ID
	uid := NewID(byte(RT_user), []byte(up.User))
	// Find user
	userRes, err := d.Get(uid, lgr)
	if err != nil {
		if err != ErrKeyNotFound {
			lgr.Error(err)
		} else {
			lgr.Error(up.User, " not found")
		}
		output <- NewMessage(req.Id, MT_fail, "")
		return
	}
	// Got a user, now cast and validate
	user, ok := userRes.(*User)
	if !ok {
		lgr.Error("cast from resource to user failed")
		output <- NewMessage(req.Id, MT_fail, "")
		return
	}

	if !user.ValidatePassword(up.Pass) {
		lgr.Error("invalid password for ", up.User)
		output <- NewMessage(req.Id, MT_fail, "")
		return
	}
	(*userid) = up.User
	lgr.Printf("%s logged in successfully", up.User)
	output <- NewMessage(req.Id, MT_ok, "")
	lgr.Time(st, "api call")
})

// SECTION:CMD

var HandlePing = WrapEndpoint(func(_ *DB, _ *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	req.Type = MT_pong
	output <- req
	lgr.Time(st, "api call")
})

var HandleList = WrapEndpoint(func(database *DB, userid *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	rt, ok := StringToResourceType(string(req.Payload))
	if !ok {
		output <- NewMessage(req.Id, MT_errinvalidformatting, "")
		return
	}
	df := database.Flow(rt)
	defer df.Stop()
	recStartTime := time.Now()
	for {
		cur, ok := df.Get()
		if !ok {
			output <- NewMessage(req.Id, MT_stop, "")
			return
		}
		enc, err := cur.MarshalString()
		if err != nil {
			lgr.Error("failed to marshal message ", err)
			output <- NewMessage(req.Id, MT_stop, "")
			return
		}
		// FIXME: ensure other methods match
		output <- NewMessage(req.Id, MessageType(req.Payload), enc)
		lgr.Time(recStartTime, "returning a resource")
		recStartTime = time.Now()
	}
	lgr.Time(st, "api call")
})

var HandleNew = WrapEndpoint(func(database *DB, userid *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	rt, ok := StringToResourceType(string(req.Payload))
	if !ok {
		lgr.Error("unable to parse resource type ", req.Payload)
		output <- NewMessage(req.Id, MT_errinvalidformatting, "")
		return
	}
	mt, ok := ResourceToMessageType(rt)
	if !ok {
		lgr.Error("unable to find matching type for resource type ", req.Payload)
		output <- NewMessage(req.Id, MT_errinvalidformatting, "")
		return
	}
	// Looks like they're cleared, so they can start sending resources
	output <- NewMessage(req.Id, MT_ok, "")

	for {
		cur := <-input
		receiveInputTime := time.Now()
		if cur == nil {
			// Then the channel has been closed
			lgr.Error("HandleNew: channel prematurely closed")
			return
		}
		switch cur.Type {
		case MT_stop:
			output <- NewMessage(req.Id, MT_ok, "")
			return
		case mt:
			res, err := UnmarshalResourceText(rt, cur.Payload)
			if err != nil {
				lgr.Error("failed to unmarshal resource of supposed type, ", string(rt), " with error ", err)
				output <- NewMessage(req.Id, MT_fail, "")
			}
			lgr.Printf("adding resource %s", res.String())
			database.Add(res)
		default:
			// Some other type -- ignore, not authorized
			output <- NewMessage(req.Id, MT_errnotauth, cur.Payload)
		}
		lgr.Time(receiveInputTime, "create new record")
	}
	lgr.Time(st, "api call")
})

var HandleGet = WrapEndpoint(func(database *DB, userid *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	want, err := ParseHex(req.Payload)
	if err != nil {
		lgr.Error("couldn't parse ID ", req.Payload)
		output <- NewMessage(req.Id, MT_errinvalidformatting, "")
		return
	}
	res, err := database.Get(want, lgr)
	if err != nil {
		if err == ErrKeyNotFound {
			lgr.Error("we don't have ID ", want.String())
			output <- NewMessage(req.Id, MT_noexists, "")
		} else {
			lgr.Error("failed to get ID ", want.String())
			output <- NewMessage(req.Id, MT_fail, "")
		}
		return
	}
	enc, err := res.MarshalString()
	if err != nil {
		lgr.Error("failed to marshal resource, res")
	}
	mt, ok := ResourceToMessageType(res.Type())
	if !ok {
		lgr.Error("failed to convert resource type to message type ", res.Type())
		output <- NewMessage(req.Id, MT_errnotimplemented, "")
		return
	}
	output <- NewMessage(req.Id, mt, enc)
	lgr.Time(st, "api call")
})

var HandleDelete = WrapEndpoint(func(database *DB, userid *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	want, err := ParseHex(req.Payload)
	if err != nil {
		lgr.Error("couldn't parse ID ", req.Payload)
		output <- NewMessage(req.Id, MT_errinvalidformatting, "")
		return
	}
	err = database.Delete(want)
	if err != nil {
		// FIXME: RESUME HERE
		// Check deleting, and return either 'ok','exists', or 'fail'
		lgr.Error("error deleting: ", err)
		// Parse error here (not found, or actual error?)
	}
	output <- NewMessage(req.Id, MT_ok, "")
	lgr.Time(st, "api call")
})

var HandleFullIndex = WrapEndpoint(func(database *DB, userid *string, req *Message, input, output chan *Message, lgr *Logger) {
	st := time.Now()
	database.FullIndex(lgr)
	output <- NewMessage(req.Id, MT_ok, "")
	lgr.Time(st, "api call")
})

func init() {
	ENDPOINTS = make(map[MessageType]Endpoint)

	ENDPOINTS[MT_login] = HandleLogin

	// CMD
	ENDPOINTS[MT_ping] = HandlePing
	ENDPOINTS[MT_list] = HandleList
	ENDPOINTS[MT_new] = HandleNew
	ENDPOINTS[MT_get] = HandleGet
	ENDPOINTS[MT_del] = HandleDelete
	ENDPOINTS[MT_fullindex] = HandleFullIndex

}
