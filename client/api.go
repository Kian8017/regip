package client

import (
	"errors"
	"regip"
	"time"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNotAuthorized  = errors.New("not authorized")
	ErrFailed         = errors.New("failed")
)

func (c *Client) Ping() bool {
	c.wg.Add(1)
	defer c.wg.Done()
	seq, inp, lg := c.Helper("ping")
	defer c.stopListen(seq)
	c.Send(regip.NewMessage(seq, regip.MT_ping, ""))
	select {
	case <-inp:
		// FIXME: make sure we got a pong
		return true
	case <-time.After(TIMEOUT):
		lg.Error("ping: timeout reached, no response from server")
		return false
	}
}

func (c *Client) List(t byte) chan regip.Resource {
	c.wg.Add(1)
	seq, inp, lg := c.Helper("list")

	mt, ok := regip.ResourceToMessageType(t)
	if !ok {
		lg.Error("Unknown resource type, ", t)
		c.wg.Done()
		return nil
	}
	// Send initial message
	c.Send(regip.NewMessage(seq, regip.MT_list, string(mt)))
	ch := make(chan regip.Resource)
	go func() {
		lg.Print("Client.List.go(): start")
		defer c.wg.Done()
		defer c.stopListen(seq)
		for {
			lg.Print("Client.List.go(): waiting for input")
			mes, ok := <-inp
			if !ok { // Channel is closed
				lg.Error("Client.List.go(): channel is closed, exiting...")
				close(ch)
				return
			}
			if mes == nil {
				lg.Error("Client.List.go(): MESSAGE IS NIL")
				continue
			}
			lg.Print("Client.List.go(): got message ", mes.Type)
			switch mes.Type {
			case regip.MT_stop:
				close(ch)
				return
			default:
				rec, err := regip.UnmarshalResourceText(t, mes.Payload)
				if err != nil {
					lg.Error("trying to unmarshal resource, ", err)
					lg.Error("\tPayload:", mes.Payload)
					continue
				}
				// Unmarshal and pass on
				ch <- rec
			}
		}
		lg.Print("Client.List.go(): end")
	}()
	return ch
}

func (c *Client) Get(i regip.ID) (regip.Resource, error) {
	c.wg.Add(1)
	defer c.wg.Done()
	seq, inp, lg := c.Helper("get")
	defer c.stopListen(seq)

	im := regip.NewMessage(seq, regip.MT_get, i.String())
	c.Send(im)
	mes := <-inp
	switch mes.Type {
	case regip.MT_errnotimplemented:
		return nil, ErrNotImplemented
	case regip.MT_errnotauth:
		return nil, ErrNotAuthorized
	case regip.MT_errinvalidformatting:
		panic("Client.Get got invalid formatting error")
		return nil, nil
	case regip.MT_fail:
		return nil, ErrFailed
	default:
		rec, err := regip.UnmarshalResourceText(i[0], mes.Payload)
		if err != nil {
			lg.Error("trying to unmarshal resource, ", err)
			lg.Error("\tPayload:", mes.Payload)
			return nil, err
		}
		return rec, nil
	}
}

func (c *Client) Delete(i regip.ID) error {
	c.wg.Add(1)
	defer c.wg.Done()
	seq, inp, lg := c.Helper("delete")
	defer c.stopListen(seq)

	im := regip.NewMessage(seq, regip.MT_del, i.String())
	c.Send(im)
	mes := <-inp
	switch mes.Type {
	case regip.MT_fail:
		return ErrFailed
	case regip.MT_ok:
		return nil
	case regip.MT_errnotauth:
		return ErrNotAuthorized
	default:
		lg.Error("unknown type of response", mes)
		return ErrNotImplemented
	}
}

func (c *Client) Add(rt byte) (chan regip.Resource, chan error, error) {
	c.wg.Add(2)
	seq, inp, lg := c.Helper("add")
	tt, ok := regip.ResourceTypeToString(rt)
	if !ok {
		return nil, nil, regip.UnknownResource
	}
	mt, ok := regip.ResourceToMessageType(rt)
	if !ok {
		panic("Resource type OK, but not message type: " + tt)
	}

	// Send start message
	im := regip.NewMessage(seq, regip.MT_new, tt)
	c.Send(im)

	res := <-inp
	if res.Type != regip.MT_ok {
		lg.Error("Initial response wasn't ok: ", res)
		e := GetError(res.Type)
		if e != nil {
			return nil, nil, e
		} else {
			lg.Error("unknown type of response ", res)
			return nil, nil, ErrNotImplemented
		}
	}

	ch := make(chan regip.Resource)
	qch := make(chan error)

	// Loop for handling messages to add
	// MESSAGE LOOP
	go func() {
		defer c.wg.Done()
		defer c.stopListen(seq)
		total := 0
		for res := range ch {
			// FIXME: RESUME here (how to handle messages from inp asynchronously? Or do I get a response back for each time I add something?)
			pay, err := res.MarshalString()
			if err != nil {
				lg.Error("Error marshaling resource: ", err)
				// Signal quit
				close(qch)
				return
			}
			// Send message
			c.Send(regip.NewMessage(seq, mt, string(pay)))
			total++
		}
		lg.Print("Sent ", total, ", now sending stop...")
		c.Send(regip.NewMessage(seq, regip.MT_stop, ""))
		lg.Print("Sent stop")
	}()

	// ERROR LOOP

	go func() {
		defer c.wg.Done()
		for {
			// Await responses of fail or not ok from the server
			select {
			case mes := <-inp:
				if mes == nil {
					lg.Print("Client.Add: got nil on error loop, exiting...")
					return
				}
				// We got something from the server
				if mes.Type != regip.MT_ok {
					e := GetError(mes.Type)
					if e == nil {
						lg.Error("unknown message type error ", mes)
						e = ErrNotImplemented
					}
					select {
					case qch <- e:
						continue
					default:
						lg.Print("send wasn't ready, exiting error loop early")
						return
					}
				}
			case <-qch:
				lg.Print("quit channel closed -- exiting error loop")
				return
			}
		}
	}()
	return ch, qch, nil
}

func (c *Client) FullIndex() error {
	c.wg.Add(1)
	defer c.wg.Done()

	seq, inp, lg := c.Helper("fullindex")
	defer c.stopListen(seq)
	lg.Print("starting fullindex")
	c.Send(regip.NewMessage(seq, regip.MT_fullindex, ""))
	mes := <-inp
	if mes.Type != regip.MT_ok {
		return GetError(mes.Type)
	}
	return nil
}
