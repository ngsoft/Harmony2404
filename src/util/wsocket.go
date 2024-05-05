package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	CONNECTION_CLOSE EventType = "conn_close"
)

type WebSocketHandler interface {
	HandleWebSocket(*WSConn)
}

type WSConn struct {
	BaseHandler
	conn  *websocket.Conn
	Route string
	Status
	ReadWriteLock
	EventListener
}

type WebSocket struct {
	Route   string
	handler WebSocketHandler
	Logger
}

func (s *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upg := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		s.Error(err.Error())
		return
	}

	h := WSConn{
		Route:  s.Route,
		conn:   conn,
		Status: STATUS_ON,
	}
	h.Initialize()
	AddEventHandler(&h, ShutdownEvent)
	defer h.Close()
	h.Info("new connection, executing handler")
	s.handler.HandleWebSocket(&h)
}

func CreateWebSocketRoute(route string, h WebSocketHandler) bool {
	if route[:1] != "/" {
		return false
	}

	ws := WebSocket{
		Route:   route,
		Logger:  NewLogger("[WebSocket]"),
		handler: h,
	}
	http.Handle(route, &ws)
	return true
}

func (h *WSConn) OnEvent(ev *Event) {
	switch ev.Type {
	case ShutdownEvent:
		h.Info("received shutdown event, closing connection")
		h.Close()
	}
}

func (h *WSConn) Close() {

	if h.On() {
		h.Status = STATUS_OFF
		h.Info("Closing connection")
		_ = h.conn.Close()
		h.DispatchEvent(CONNECTION_CLOSE)
		h.Info("Connection closed")
	}

}

func (h *WSConn) ReadMessage() (string, bool) {
	var (
		mt    int
		msg   string
		input []byte
		err   = fmt.Errorf("connection status is off or undefined")
	)

	if h.On() {
		h.ReadLock.Lock()
		defer h.ReadLock.Unlock()
		if mt, input, err = h.conn.ReadMessage(); err == nil {
			if mt == websocket.BinaryMessage {
				err = fmt.Errorf("binary message received")
			} else {
				msg = strings.TrimRight(string(input), "\r\n")
			}
		}
	}

	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}

	return msg, err == nil
}

func (h *WSConn) WriteMessage(message string) bool {
	var err = fmt.Errorf("connection status is off or undefined")
	if h.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.conn.WriteMessage(websocket.TextMessage, []byte(message))
	}

	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}
	return err == nil
}

func (h *WSConn) WriteJson(v interface{}) bool {
	var err = fmt.Errorf("connection status is off or undefined")
	if h.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.conn.WriteJSON(v)
	}
	if err != nil && h.On() {
		h.Error(err.Error())
		h.Close()
	}

	return err == nil
}
