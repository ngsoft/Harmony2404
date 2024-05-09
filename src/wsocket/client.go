package wsocket

import (
	"encoding/json"
	"flirc/util"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	util.BaseHandler
	util.Status
	Conn        *websocket.Conn
	WebSocket   *WebSocket
	CurrentRoom *Room
	pong        func()
	ok          chan bool
	running     bool
	Handler     Handler
	util.ReadWriteLock
}

func (c *Client) OnMessage(d Direction, t Type, v ...interface{}) {
	var (
		ok bool
	)

	c.Info("client event(direction=>`%s`, type=>`%s`, args=>(%v))", d.String(), t, v)
	// handle reserved event
	if d == In && t == JoinRoom {
		room := ""
		if len(v) > 0 {
			if room, ok = v[0].(string); ok {
				ok = c.WebSocket.SwitchRoom(c, room)
			}
		}

		if !ok {
			c.SendEvent(Error, "invalid room "+room)
		}
		return
	}

	// external handler
	if c.Handler != nil {
		c.Handler.OnMessage(c, d, t, v...)
	} else if c.CurrentRoom != nil {
		// send event to room
		c.CurrentRoom.OnMessage(c, d, t, v...)
	}
}

func (c *Client) Run() {
	c.Info("start")
	if !c.running {
		c.running = true
		if c.WebSocket.Ping {
			go c.pingHandler()
		}
		go c.reader()
		go c.WebSocket.DispatchEvent(Open, c)
	}

}

func (c *Client) Close() {
	if c.On() {
		c.Status++
		close(c.ok)
		c.WebSocket.RemoveClient(c)
		c.pong()
		_ = c.Conn.Close()
		go c.WebSocket.DispatchEvent(Close, c)
	}
}
func (c *Client) reader() {

	var (
		msg   string
		ok    bool
		lines []string
		t     Type
		tstr  string
		err   error
		args  Payload
		v     []interface{}
	)
	defer c.Close()

	for {
		if msg, ok = c.ReadMessage(); !ok {
			return
		}

		// clean up msg (windows clients)
		msg = strings.ReplaceAll(msg, "\r", "")
		// lines
		lines = strings.Split(msg, "\n")
		// is event
		for _, line := range lines {
			// empty line
			if len(line) == 0 {
				continue
			}

			t = Message
			v = []interface{}{line}
			if line[:1] == "[" && line[len(line)-1:] == "]" {
				args = Payload{}
				err = json.Unmarshal([]byte(line), &args)
				if err != nil || len(args) == 0 {
					c.Error("cannot unmarshal=>%v", line)
					return
				}

				c.Info("args: %v", args)
				if tstr, ok = args[0].(string); ok {
					t = Type(tstr)
					v = make([]interface{}, 0)
					for i := 1; i < len(args); i++ {
						v = append(v, args[i])
					}
				} else {
					c.Error("received invalid data=>%v", line)
					return
				}

			}

			// call message handler
			go c.OnMessage(In, t, v...)

		}

	}

}
func (c *Client) pingHandler() {
	var (
		interval    = time.Duration(c.WebSocket.PingInterval) * time.Second
		pongTimeout = (interval * 110) / 100
		ticker      = time.NewTicker(interval)
		msg         = "{\"op\":\"ping\"}"
		ok          bool
		fn          = func() {
			c.Error("pong timeout, disconnecting client")
			c.ok <- false
		}
	)

	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case valid, open := <-c.ok:
			if !open || !valid {
				return
			}
		case <-ticker.C:
			if !c.On() {
				return
			}
			if ok = c.sendPing(msg); ok {
				c.Info("ping sent to client")
				c.pong = util.SetTimeout(fn, pongTimeout)
			}
		}
	}

}
func (c *Client) SendEvent(t Type, v ...interface{}) bool {
	var (
		err   error
		input []byte
		p     = Payload{
			t.String(),
		}
	)

	for _, arg := range v {
		p = append(p, arg)
	}

	input, err = json.Marshal(p)
	if err != nil {
		c.Error("cannot marshal=>%v", p)
		return false
	}
	if c.SendMessage(string(input)) {
		go c.OnMessage(Out, t, v...)
		return true
	}
	return false
}

func (c *Client) SendMessage(message string) bool {
	var err = fmt.Errorf("connection status is off")
	if c.On() {
		c.WriteLock.Lock()
		defer c.WriteLock.Unlock()
		err = c.Conn.WriteMessage(websocket.TextMessage, []byte(strings.TrimRight(message, "\r\n")+"\n"))
	}

	if err != nil && c.On() {
		c.Error(err.Error())
		c.Close()
	}
	return err == nil
}

func (c *Client) ReadMessage() (string, bool) {
	var (
		loop  = true
		mt    int
		msg   string
		input []byte
		err   = fmt.Errorf("connection status is off")
	)

	if c.On() {
		c.ReadLock.Lock()
		defer c.ReadLock.Unlock()

		conn := c.Conn

		for loop {
			loop = false
			if mt, input, err = conn.ReadMessage(); err == nil {
				switch mt {
				case websocket.BinaryMessage:
					err = fmt.Errorf("binary message received")
				case websocket.TextMessage:
					msg = string(input)
				case websocket.PingMessage:
					c.Info("received ping, sending pong")
					err = fmt.Errorf("cannot send pong message")
					if ok := c.sendPong(string(input)); ok {
						// wait for real message
						loop = true
						err = nil
					}

				case websocket.PongMessage:
					c.Info("received pong message")
					c.pong()
					loop = true
				case websocket.CloseMessage:
					err = fmt.Errorf("received close message")
				}
			}
		}

	}

	if err != nil && c.On() {
		c.Error(err.Error())
		c.Close()
	}

	return msg, err == nil
}

func (c *Client) sendPing(message string) bool {
	var err = fmt.Errorf("connection status is off")
	if c.On() {
		c.WriteLock.Lock()
		defer c.WriteLock.Unlock()
		err = c.Conn.WriteMessage(websocket.PingMessage, []byte(strings.TrimRight(message, "\r\n")+"\n"))
	}
	if err != nil && c.On() {
		c.Error(err.Error())
		c.Close()
	}
	return err == nil
}

func (c *Client) sendPong(message string) bool {
	var err = fmt.Errorf("connection status is off")
	if c.On() {
		c.WriteLock.Lock()
		defer c.WriteLock.Unlock()
		message = strings.TrimRight(message, "\n") + "\n"
		err = c.Conn.WriteMessage(websocket.PongMessage, []byte(strings.TrimRight(message, "\r\n")+"\n"))
	}
	if err != nil && c.On() {
		c.Error(err.Error())
		c.Close()
	}
	return err == nil
}
