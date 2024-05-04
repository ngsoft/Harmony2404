package main

import (
	"flirc/util"
	"strings"
)

func (h *FlircHandler) HandleEvent(ev *util.Event) {
	switch ev.Type {
	case util.CONNECTION_CLOSE:
		h.Info("Connection closed ")
	case INPUT_EVENT:
		h.Log(h.Message.Message)
	}

}

func (h *Handler) HandleEvent(ev *util.Event) {
	switch ev.Type {
	case INPUT_EVENT:
	case util.SHUTDOWN_EVENT:
	}
}

func (h *ConnHandler) HandleWebSocket(ws *util.WSConn) {

	h.Info("new connection to websocket")
	var (
		msg     string
		success bool
	)
	for {

		msg, success = ws.ReadMessage()
		if !success {
			if ws.On() {
				util.Shutdown(1)
			}
			break
		}

		ws.Log(msg)

	}

}
func (h *FlircHandler) HandleUnixSocket(s *util.UnixSocket) {

	var (
		msg     string
		success bool
		data    []string
	)

	h.Info("connected to socket %v", s.Path)

	for {
		msg, success = s.ReadMessage()
		if !success {
			if s.Status.On() {
				h.Error("disconnected from socket %v", s.Path)
				util.Shutdown(1)
			}
			return
		}

		// decode packet
		data = strings.Split(msg, " ")
		if len(data) != 4 {
			h.Warn("wrong packet received: `%s`", msg)
			continue
		}

		if h.remote != data[3] {
			continue
		}

		h.Message = SocketMessage{
			Message: msg,
			Code:    data[0],
			Repeat:  util.HexToInt(data[1]),
			Key:     data[2],
			Remote:  data[3],
		}

		h.DispatchEvent(INPUT_EVENT)

	}

}

// func (h *Handler) listenToWebSocket(ws *util.WebSocket) {

// 	h.Info("new connection to websocket")
// 	var (
// 		msg     string
// 		success bool
// 	)
// 	for {

// 		msg, success = ws.ReadMessage()
// 		if !success {
// 			h.Error("disconnected from websocket")
// 			util.GlobalEvents.Shutdown(1)
// 		}

// 		h.Log(msg)

// 	}

// }

// func HandleWebookRoute(w http.ResponseWriter, r *http.Request) {

// 	var base util.BaseHandler = util.NewBaseHandler()

// 	ws, err := util.NewWebsocket(w, r)
// 	if err != nil {
// 		base.Error(err.Error())
// 		util.GlobalEvents.Shutdown(1)
// 	}

// 	h := Handler{BaseHandler: base, WebSocket: ws}

// 	ws.HandleConnection(h.listenToWebSocket)

// }
