package regip

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"sync"
)

const DEFAULT_LISTEN_ADDR = ":8089"

type Server struct {
	DB     *DB
	API    *API
	Logger *Logger

	Addr     string
	Server   *http.Server
	ServeMux *http.ServeMux
	upgrader websocket.Upgrader
}

func NewServer(dbLoc, listen string, log bool) (*Server, error) {
	d, err := NewDB(dbLoc)
	if err != nil {
		return nil, err
	}
	var s Server
	s.DB = d
	s.Addr = listen
	s.API = NewAPI(s.DB)

	// Logger
	// FIXME check log-bool here
	s.Logger = NewLogger(os.Stdout).Tag("server", CLR_server)

	// ServeMux
	s.ServeMux = http.NewServeMux()
	// Add serveMux routes here
	s.ServeMux.Handle("/", http.FileServer(FS(false))) // false = don't use local
	s.ServeMux.HandleFunc("/ws", s.HandleWS)

	// Server
	s.Server = &http.Server{}
	// Assign addr + serveMux
	s.Server.Addr = s.Addr
	s.Server.Handler = s.ServeMux

	// Websocket Upgrader
	s.upgrader = websocket.Upgrader{CheckOrigin: s.ValidateOrigin}
	s.Logger.Print("started")
	return &s, nil
}

func (s *Server) Run() error {
	s.Logger.Print("server running on", s.Addr)
	return s.Server.ListenAndServe()
}

func (s *Server) Close() {
	s.DB.Close()
}

const MAX_SESSION_BUFFER_SIZE = 10 // FIXME should be able to do this with 0

type Chain struct {
	Data chan *Message
	IsOk bool
	Mut  sync.Mutex
}

func (c *Chain) OK() bool {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	return c.IsOk
}

func (c *Chain) Stop() {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	c.IsOk = false
	close(c.Data)
}

func (c *Chain) Push(m *Message) bool {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	if !c.IsOk {
		return false
	}
	c.Data <- m
	return true
}

func NewChain() *Chain {
	var c Chain
	c.IsOk = true
	c.Data = make(chan *Message)
	return &c
}

type Session struct {
	conn     *websocket.Conn
	server   *Server
	logger   *Logger
	user     string
	chains   map[int]*Chain
	chainMut sync.RWMutex
	outgoing chan *Message
	quit     chan struct{}
}

func (s *Server) NewSession(c *websocket.Conn) *Session {
	var ns Session
	var m sync.RWMutex
	ns.conn = c
	ns.server = s

	// Logger
	sessionTag := "session " + c.RemoteAddr().String()
	ns.logger = s.Logger.Tag(sessionTag, CLR_session)
	ns.logger.Print("new")

	ns.chainMut = m

	ns.chains = make(map[int]*Chain)
	ns.outgoing = make(chan *Message, MAX_SESSION_BUFFER_SIZE)
	ns.quit = make(chan struct{})

	go ns.ReadLoop()
	go ns.WriteLoop()
	return &ns
}

func (s *Session) Close() {
	// Close all the chains

	// We already close the channels in ChainLoop
	/*
		for _, cc := range s.chains {
			cc.Stop()
		}
	*/

	close(s.quit)
	close(s.outgoing)
	s.conn.Close()
	s.logger.Print("close")
}

func (s *Session) Run() {
	<-s.quit
}

func (s *Session) ReadLoop() {
	lgr := s.logger.Tag("readloop", CLR_readloop)
	lgr.Print("new")
	for {
		_, raw, err := s.conn.ReadMessage()
		if err != nil {
			s.Close()
			return
		}

		// Get the message
		mess, err := UnmarshalMessage(raw)
		if err != nil {
			lgr.Error(err)
			continue
		}

		s.ChainLoop(mess)
	}
	lgr.Print("end")
}

func (s *Session) ChainLoop(m *Message) {
	// FIXME: these two functions, ReadLoop and ChainLoop, are unnecessarily separated
	lgTag := fmt.Sprintf("chainloop %d", m.Id)
	lgr := s.logger.Tag(lgTag, CLR_chainloop)

	// DEBUG
	lgr.RawMessage(m)
	var ap bool

	// Lock to see if we have a chain for that ID
	// NOT RLOCK (because we need to immediately claim it under the same mutex if it's free)
	s.chainMut.Lock()

	// Check for existence, or make here
	cur, ok := s.chains[m.Id]
	if ok { // we already have a chain
		ap = true
	} else {
		s.chains[m.Id] = NewChain()
	}

	s.chainMut.Unlock()
	// FIXME: RESUME HERE

	if ap {
		select {
		case <-s.quit:
		default:
			cur.Push(m)
		}
		return
	}

	// This needs to be in a separate function
	go func() {
		lgr.Print("new chain")
		defer lgr.Print("end chain")
		// FIXME: DEBUG: we should probably pass the whole chain, not just a channel
		quit := s.server.API.Handle(m, s.chains[m.Id].Data, s.outgoing, &s.user, lgr)
		select {
		case <-quit:
			s.chains[m.Id].Stop()
			s.chainMut.Lock()
			delete(s.chains, m.Id)
			s.chainMut.Unlock()
		case <-s.quit:
			s.chains[m.Id].Stop()
			s.chainMut.Lock()
			delete(s.chains, m.Id)
			s.chainMut.Unlock()
		}
	}()
}

func (s *Session) WriteLoop() {
	lgr := s.logger.Tag("writeloop", CLR_writeloop)
	lgr.Print("new")
	for outgoingMessage := range s.outgoing {
		// DEBUG
		lgr.RawMessage(outgoingMessage)
		encoded, err := outgoingMessage.Marshal()
		if err != nil {
			lgr.Error("marshaling outgoing message", err)
			continue
		}
		err = s.conn.WriteMessage(websocket.TextMessage, encoded)
		if err != nil {
			lgr.Error("writing message", err)
			s.Close()
		}
	}
	lgr.Print("end")
}

func (s *Server) ValidateOrigin(r *http.Request) bool {
	// FIXME: Implement
	return true
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte("Invalid request sent to websocket URL"))
		return
	}
	sess := s.NewSession(c)
	sess.Run()
}
