package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func (h *Handler) Start() {

	http.Handle(h.Route, h)
	h.info("Starting web socket on port %v", h.Port)
	h.error("%v", http.ListenAndServe(fmt.Sprintf(":%d", h.Port), nil))
	h.Shutdown(1)
}

func (h *Handler) handleEvent(ev string) error {
	var err error = fmt.Errorf("invalid event %v", ev)
	switch ev {
	case "start":
		h.info("received start event")
		err = nil
		h.handleUnixSocket()
	}
	return err
}
func (h *Handler) handleWebSocket() {
	var (
		msg string
		err error = nil
	)
	for {
		msg, err = h.readMessage()
		if err != nil {
			h.error("cannot read ws message: %v", err)
			return
		}
		go func() {
			if err = h.handleEvent(msg); err != nil {
				h.error("cannot handle message: %v", err)
				return
			}
		}()

	}

	// h.handleEvent("start")
}

func (h *Handler) handleUnixSocket() error {
	var (
		mess string
		cnt  int
		err  error = nil
	)
	h.info("connecting to [%v]", h.Socket)

	c, err := net.Dial("unix", h.Socket)
	if err != nil {
		return err
	}

	h.info("connection open to [%v]", h.Socket)

	h.UnixSocket = UnixSocket{c, true}
	defer c.Close()
	for {
		buf := make([]byte, 1024)
		cnt, err = c.Read(buf)
		if err != nil {
			if h.UnixSocket.Open {
				h.error("handleConn %v", err)
			}
			return err
		}
		mess = strings.TrimRight(string(buf[0:cnt]), "\r\n")

		if len(mess) == 0 {
			h.error("empty message")
			continue
		}
		if data, err := h.processMessage(mess); err == nil {
			go h.handleMessage(data)
		}
	}
}

func (h *Handler) processMessage(message string) (SocketMessage, error) {
	var result SocketMessage
	data := strings.Split(message, " ")
	if len(data) != 4 {
		err := fmt.Errorf("not processing message invalid message length %v instead of 4, message => `%v`", len(data), message)
		h.error(err.Error())
		return result, err
	}
	result = SocketMessage{
		Message: message,
		Code:    data[0],
		Repeat:  hexToInt(data[1]),
		Key:     data[2],
		Remote:  data[3],
	}
	return result, nil
}

func (h *Handler) handleMessage(mess SocketMessage) error {

	if mess.Remote != h.Remote {
		h.info("received unmanaged packet from remote `%v`, managed => `%v`", mess.Remote, h.Remote)
		return nil
	}
	h.log(mess.Message)
	return h.writeJson([]interface{}{
		fmt.Sprintf("0x%s", mess.Code),
		mess.Key,
		mess.Repeat,
	})
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsh := &h.WebSocket
	wsh.Writer = w
	conn, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.error("Cannot create ws:// connection %v", err)
		return
	}
	wsh.Conn = conn
	wsh.Open = true
	defer conn.Close()
	h.handleWebSocket()
}

func (h *Handler) Shutdown(code int) {
	h.info("Received SIGTERM signal, closing connections")
	if h.UnixSocket.Open {
		h.UnixSocket.Open = false
		h.UnixSocket.Conn.Close()
		h.info("closing websocket connection")
	}

	if h.WebSocket.Open {
		h.WebSocket.Open = false
		h.WebSocket.Conn.Close()
		h.info("connection close to [%v]", h.Socket)
	}

	os.Exit(code)
}
