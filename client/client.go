package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"regip"
	"sync"
	"time"
)

const TIMEOUT = 3 * time.Second

type Client struct {
	conn   *websocket.Conn
	logger *regip.Logger
	wg     *sync.WaitGroup

	// Sequence IDs
	sequence int
	seqMut   sync.Mutex

	// Listens
	listens   map[int]chan *regip.Message
	listenMut sync.Mutex

	// Outgoing
	outgoing chan *regip.Message
	//quit     chan struct{}
}

func NewClient(addr string, lgr *regip.Logger) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	var c Client
	c.conn = conn
	c.logger = lgr

	var wg sync.WaitGroup
	c.wg = &wg

	c.listens = make(map[int]chan *regip.Message)

	c.outgoing = make(chan *regip.Message)
	//c.quit = make(chan struct{})

	go c.ReadLoop()
	go c.WriteLoop()

	return &c, nil
}

func (c *Client) Wait() {
	c.logger.Print("waiting for waitgroup")
	c.wg.Wait()
	close(c.outgoing)

}

func (c *Client) ReadLoop() {
	lgr := c.logger.Tag("readloop", regip.CLR_readloop)
	lgr.Print("start")
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			//close(c.quit)
			lgr.Error("end (read error or close)")
			return
		}
		// Get the message
		mess, err := regip.UnmarshalMessage(raw)
		if err != nil {
			lgr.Error(err)
			continue
		}

		// DEBUG
		lgr.RawMessage(mess)

		// See if we have a chain for that...

		// DEBUG
		// lgr.Print("Checking for chain")

		c.listenMut.Lock()
		ch, ok := c.listens[mess.Id]
		c.listenMut.Unlock()

		// DEBUG
		// lgr.Print("Done checking for chain")

		if !ok { // No one is listening for this ID currently
			lgr.Print("no chain, continuing")
			continue
		} else { // Someone's listening...
			lgr.Print("adding to existing chain...")
			ch <- mess
		}
	}
	lgr.Print("end")
}

func (c *Client) WriteLoop() {
	lgr := c.logger.Tag("writeloop", regip.CLR_writeloop)
	lgr.Print("start")
	for mes := range c.outgoing {
		// send message
		encoded, err := mes.Marshal()
		if err != nil {
			lgr.Error("marshaling outgoing message", err)
			continue
		}

		err = c.conn.WriteMessage(websocket.TextMessage, encoded)
		if err != nil {
			lgr.Error("writing message", err)
			//close(c.quit)
			return
		}
		// DEBUG
		lgr.RawMessage(mes)
	}
	lgr.Print("end")
}

func (c *Client) getSeq() int {
	c.seqMut.Lock()
	defer c.seqMut.Unlock()
	c.sequence++
	return c.sequence
}

func (c *Client) listen(i int) (chan *regip.Message, bool) {

	// FIXME: DEBUG
	// ft := c.logger.Tag("Client.listen", regip.CLR_mt)
	// ft.Print("start")
	// defer ft.Print("end")

	c.listenMut.Lock()
	defer c.listenMut.Unlock()
	_, ok := c.listens[i]
	if ok { // Someone's already listening -- abort
		return nil, false
	}
	ch := make(chan *regip.Message)
	c.listens[i] = ch
	return ch, true
}

func (c *Client) stopListen(i int) {
	c.listenMut.Lock()
	defer c.listenMut.Unlock()
	ch, ok := c.listens[i]
	if ok { // Only close the channel if it exists
		close(ch)
	}
	delete(c.listens, i)
}

func (c *Client) Helper(name string) (int, chan *regip.Message, *regip.Logger) {

	// FIXME: DEBUG
	// ft := c.logger.Tag("Client.Helper", regip.CLR_mt)
	// ft.Print("start")
	// defer ft.Print("end")

	seq := c.getSeq()
	lgTag := fmt.Sprintf("%s %d", name, seq)
	lg := c.logger.Tag(lgTag, regip.CLR_api)
	ch, ok := c.listen(seq)
	if !ok {
		lg.Error("sequence number ", seq, "already taken, trying again")
		return c.Helper(name)
	}
	return seq, ch, lg
}

func (c *Client) Send(m *regip.Message) {
	c.outgoing <- m
}

func GetError(mt regip.MessageType) error {
	switch mt {
	case regip.MT_fail:
		return ErrFailed
	case regip.MT_errnotauth:
		return ErrNotAuthorized
	case regip.MT_errinvalidformatting:
		// FIXME: DEBUG (this should never happen)
		panic(fmt.Sprint("client got invalid formatting error", mt))
		return nil
	case regip.MT_errnotimplemented:
		return ErrNotImplemented
	default:
		// FIXME: DEBUG
		panic(fmt.Sprint("Client.GetError -- don't have case for ", mt))
		return nil
	}
}
