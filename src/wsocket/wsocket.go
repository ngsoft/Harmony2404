package wsocket

import (
	"flirc/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	DefaultPingInterval = 60
	DefaultPingEnabled  = false
	DefaultRoute        = "/ws"
	DefaultPort         = 9023
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Config struct {
	Route        string
	Ping         bool
	PingInterval int
	Port         int
}

type WebSocket struct {
	util.Logger
	Config
	util.EventListener
	Clients map[*Client]bool
	Rooms   map[string]*Room
	init    bool
	running bool
}

func (s *WebSocket) __init() {
	if !s.init {
		s.init = true
		if s.Config.Route == "" {
			s.Config.Route = DefaultRoute
		}
		if s.Config.PingInterval == 0 {
			s.Config.PingInterval = DefaultPingInterval
		}
		if s.Config.Port == 0 {
			s.Config.Port = DefaultPort
		}
		s.Clients = make(map[*Client]bool)
		s.Rooms = make(map[string]*Room)
		s.SetLoggerPrefix("[ws]")
		s.AddEventHandler(s)
		http.Handle(s.Route, s)
	}
}

func (s *WebSocket) addr() string {
	return ":" + strconv.Itoa(s.Port)
}
func (s *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.__init()

	conn, err := Upgrader.Upgrade(w, r, nil)

	if err != nil {
		s.Error(err.Error())
		return
	}
	c := s.AddClient(conn)
	go c.Run()
}

func (s *WebSocket) OnEvent(ev *util.Event) {

	var (
		ok     bool
		client *Client
	)

	if len(ev.Params) == 0 {
		s.Error("event[%s] no client defined", ev.Type.String())
		return
	}
	if client, ok = ev.Params[0].(*Client); !ok {
		s.Error("event[%s] no client defined", ev.Type.String())
		return
	}
	switch ev.Type {
	case Open:
		s.Info("new incoming connection from client[%v]", client.Uid)

	case Close:
		s.Info("closed connection for client[%v]", client.Uid)
		s.RemoveClient(client)
	}

}

func (s *WebSocket) AddClient(c *websocket.Conn) *Client {
	s.__init()
	var client = Client{
		WebSocket: s,
		Status:    util.NewStatus(),
		Conn:      c,
		pong:      func() {},
		ok:        make(chan bool),
	}
	go func() {
		client.ok <- true
	}()
	client.Initialize()
	s.Clients[&client] = true
	return &client
}

func (s *WebSocket) RemoveClient(c *Client) {
	s.__init()
	if c.CurrentRoom != nil {
		c.CurrentRoom.RemoveClient(c)
	}
	delete(s.Clients, c)
}

func (s *WebSocket) HasRoom(name string) bool {
	s.__init()
	_, ok := s.Rooms[name]
	return ok
}

func (s *WebSocket) AddRoom(name string) *Room {
	s.__init()
	if !s.HasRoom(name) {
		s.Rooms[name] = &Room{
			Name:      name,
			Clients:   make(map[*Client]bool),
			WebSocket: s,
		}
	}

	return s.Rooms[name]
}

func (s *WebSocket) AddRooms(names []string) []*Room {
	result := make([]*Room, 0)
	for _, name := range names {
		result = append(result, s.AddRoom(name))
	}
	return result
}

func (s *WebSocket) SwitchRoom(c *Client, name string) bool {
	s.__init()
	if !s.HasRoom(name) {
		return false
	}
	room := s.Rooms[name]
	room.AddClient(c)
	return true
}

func (s *WebSocket) Run() {
	s.__init()
	var err = fmt.Errorf("websocket already running")
	if !s.running {
		s.running = true
		s.Info("Starting websocket on ws://127.0.0.1%s%s", s.addr(), s.Route)
		err = http.ListenAndServe(s.addr(), nil)
	}

	if err != nil {
		s.Error(err.Error())
	}
}
