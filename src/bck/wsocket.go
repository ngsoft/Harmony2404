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
type WebSocket struct {
	Route   string
	handler WebSocketHandler
	Logger
}
type WSConn struct {
	conn *websocket.Conn
	Status
	ReadWriteLock
	Logger
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

func (h *WebSocket) getUpgrader() *websocket.Upgrader {

	if !h.init {
		upg := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}

	}

}

func (h *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upg := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		h.Error(err.Error())
		return
	}

}

func NewWebsocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	upg := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	ws := WebSocket{
		Status: STATUS_OPEN,
		conn:   conn,
	}

	GlobalEvents.AttachEvent(func(e *Event) {
		if ws.Status == STATUS_OPEN {
			ws.Status = STATUS_CLOSE
			_ = ws.conn.Close()
		}
	}, SHUTDOWN_EVENT)

	return &ws, nil
}

func (h *WebSocket) HandleConnection(handler func(*WebSocket)) {
	if h.Status == STATUS_OPEN {
		defer h.conn.Close()
		handler(h)
	}
	h.Status = STATUS_CLOSE
}

func (h *WebSocket) ReadMessage() (string, bool) {
	var (
		mt    int
		msg   string
		input []byte
		err   error = fmt.Errorf("")
	)

	if h.Status == STATUS_OPEN {
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

	return msg, err == nil
}

func (h *WebSocket) WriteMessage(message string) bool {

	if h.Status == STATUS_OPEN {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		return h.conn.WriteMessage(websocket.TextMessage, []byte(message)) == nil
	}
	return false
}

func (h *WebSocket) WriteJson(v interface{}) bool {
	if h.Status == STATUS_OPEN {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()
		return h.conn.WriteJSON(v) == nil
	}
	return false
}
