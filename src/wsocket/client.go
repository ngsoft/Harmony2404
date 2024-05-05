package wsocket

import (
	"flirc/util"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/imbhargav5/noop"
)

const (
	ErrorEvent EventType = "error"
	pingEvent  EventType = "ping"
)

type ClientHandler interface {
	OnMessage(*Client, string)
	OnMessageEvent(*Client, *Event)
	OnOpen(*Client)
	OnClose(*Client)
}

type clientconfig struct {
	pingInterval time.Duration
	pongTimeout  time.Duration
	mustPing     bool
}

type Client struct {
	Handler ClientHandler
	util.BaseHandler
	util.Status
	Conn *websocket.Conn
	util.ReadWriteLock
	incoming chan *Event
	outgoing chan *Event
	pong     func()
	ok       chan bool
	running  bool
	clientconfig
	Locked bool
}

func NewClient(r *Router, c *websocket.Conn) *Client {

	cl := Client{
		incoming: make(chan *Event),
		outgoing: make(chan *Event),
		Handler:  r.RouterConfig.ClientHandler,
		Conn:     c,
		Status:   util.NewStatus(),
		pong:     noop.Noop,
		ok:       make(chan bool),
		clientconfig: clientconfig{
			mustPing:     r.RouterConfig.PingEnabled,
			pingInterval: r.RouterConfig.PingInterval,
			pongTimeout:  (r.RouterConfig.PingInterval * 110) / 100,
		},
	}
	cl.Initialize()
	cl.ok <- true
	return &cl
}
func (c *Client) Run() {
	if !c.running {
		c.running = true
		if c.mustPing {
			go c.pingHandler()
		}

		go c.messageHandler()
		c.Handler.OnOpen(c)
	}

}

func (c *Client) SetHandler(h ClientHandler) {
	c.Handler = h
}

func (h *Client) Close() {
	if h.On() {
		h.Status++
		h.Info("Closing connection")
		_ = h.Conn.Close()
		close(h.ok)
		close(h.incoming)
		close(h.outgoing)
		h.Info("Connection closed")
		h.Handler.OnClose(h)
	}
}

// handle event messages listeners
func (h *Client) messageHandler() {
	go func() {
		for {
			select {
			case ev, ok := <-h.incoming:
				if !ok {
					return
				}
				if !h.Locked {
					h.Handler.OnMessageEvent(h, ev)
				}
			case ev, ok := <-h.outgoing:
				if !ok || !h.SendMessage(ev.String()) {
					return
				}
				h.Handler.OnMessageEvent(h, ev)
			}
		}
	}()
	var (
		msg string
		ok  bool
		e   Event
	)
	defer h.Close()

	for {
		if msg, ok = h.ReadMessage(); !ok {
			return
		}
		// Event detection
		if msg[:1] == "[" && msg[len(msg)-1:] == "]" { //check if json list
			if e, ok = NewEventFromString(msg); ok {
				// send to incoming
				h.incoming <- &e
				continue
			}
		}
		// not an event: call onMessage
		h.Handler.OnMessage(h, msg)
		break
	}

}

func (h *Client) pingHandler() {

	var (
		ticker = time.NewTicker(h.pingInterval)
		msg    = NewEvent(pingEvent, "").String()
		ok     bool
		fn     = func() {
			h.Error("pong timeout")
			h.ok <- false
		}
	)

	defer func() {
		ticker.Stop()
		h.Close()
	}()

	for {
		select {
		case val, open := <-h.ok:
			if !open {
				return
			}
			if !val {
				h.Error("ping not responded")
				return
			}
		case <-ticker.C:
			if ok = h.sendPing(msg); ok {
				h.pong = util.SetTimeout(fn, h.pongTimeout)
			}
			if !ok {
				h.Error("cannot ping")
				return
			}
		}
	}

}

func (h *Client) ReadMessage() (string, bool) {
	var (
		loop  = true
		mt    int
		msg   string
		input []byte
		err   = fmt.Errorf("connection status is off")
	)

	if h.On() {
		h.ReadLock.Lock()
		defer h.ReadLock.Unlock()

		for loop {
			loop = false
			if mt, input, err = h.Conn.ReadMessage(); err == nil {
				switch mt {
				case websocket.BinaryMessage:
					err = fmt.Errorf("binary message received")
				case websocket.TextMessage:
					msg = strings.TrimRight(string(input), " \r\n")
				case websocket.PingMessage:
					h.Info("received ping, sending pong")
					err = fmt.Errorf("cannot send pong message")
					if ok := h.sendPong(strings.TrimRight(string(input), " \r\n")); ok {
						// wait for real message
						loop = true
						err = nil
					}

				case websocket.PongMessage:
					h.Info("received pong message")
					h.pong()
					loop = true
				case websocket.CloseMessage:
					err = fmt.Errorf("received close message")
				}
			}
		}

	}

	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}

	return msg, err == nil
}

func (h *Client) sendPing(message string) bool {
	var err = fmt.Errorf("connection status is off")
	if h.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.Conn.WriteMessage(websocket.PingMessage, []byte(message))
	}
	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}
	return err == nil
}

func (h *Client) sendPong(message string) bool {
	var err = fmt.Errorf("connection status is off")
	if h.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.Conn.WriteMessage(websocket.PongMessage, []byte(message))
	}
	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}
	return err == nil
}

func (h *Client) SendMessage(message string) bool {
	var err = fmt.Errorf("connection status is off")
	if h.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}

	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}
	return err == nil
}

func (h *Client) SendEvent(e *Event) {
	e.Direction = Outgoing
	h.outgoing <- e
}
