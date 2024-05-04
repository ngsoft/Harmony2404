package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type WebSocketHandler interface {
	HandleWebSocket(*WSConn)
}

type WSConn struct {
	conn *websocket.Conn
	Status
	ReadWriteLock
	Logger
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
		conn:   conn,
		Status: STATUS_ON,
		Logger: NewLogger(fmt.Sprintf("[WebSocket|%s]", s.Route)),
	}
	GlobalEvents.AddEventHandler(&h, SHUTDOWN_EVENT)
	defer h.Close()
	h.Info("new connection, executing handler")
	s.handler.HandleWebSocket(&h)
}

func CreateWebSocketRoute(route string, h WebSocketHandler) bool {
	if route[:1] != "/" {
		return false
	}
	http.Handle(route, &WebSocket{
		Route:   route,
		Logger:  NewLogger("[WebSocket]"),
		handler: h,
	})
	return true
}

func (h *WSConn) HandleEvent(ev *Event) {
	switch ev.Type {
	case SHUTDOWN_EVENT:
		h.Close()
	}
}

func (h *WSConn) Close() {

	if h.Status.On() {
		h.Info("Closing Connection")
		_ = h.conn.Close()
	}
	h.Status = STATUS_OFF
}

func (h *WSConn) ReadMessage() (string, bool) {
	var (
		mt    int
		msg   string
		input []byte
		err   error = fmt.Errorf("connection status is off or undefined")
	)

	if h.Status.On() {
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

	if err != nil {
		h.Error(err.Error())
		h.Close()
	}

	return msg, err == nil
}

func (h *WSConn) WriteMessage(message string) bool {
	var err error = fmt.Errorf("connection status is off or undefined")
	if h.Status.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.conn.WriteMessage(websocket.TextMessage, []byte(message))
	}

	if err != nil {
		h.Error(err.Error())
		h.Close()
	}

	return err == nil
}

func (h *WSConn) WriteJson(v interface{}) bool {
	var err error = fmt.Errorf("connection status is off or undefined")
	if h.Status.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		err = h.conn.WriteJSON(v)
	}
	if err != nil {
		h.Error(err.Error())
		h.Close()
	}

	return err == nil
}
