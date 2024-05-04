package main

import (
	"flirc/util"
	"net/http"
	"strings"
)

func (h *FlircHandler) handleConn(s *util.UnixSocket) {

	var (
		msg     string
		success bool
		data    []string
	)

	h.Info("connected to socket %v", s.Path)

	for {
		msg, success = s.ReadMessage()
		if !success {
			h.Error("disconnected from socket %v", s.Path)
			util.GlobalEvents.Shutdown(1)
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

		h.TriggerEvent(INPUT_EVENT)

	}

}

func (h *Handler) listenToWebSocket(ws *util.WebSocket) {

	h.Info("new connection to websocket")
	var (
		msg     string
		success bool
	)
	for {

		msg, success = ws.ReadMessage()
		if !success {
			h.Error("disconnected from websocket")
			util.GlobalEvents.Shutdown(1)
		}

		h.Log(msg)

	}

}

func HandleWebookRoute(w http.ResponseWriter, r *http.Request) {

	var base util.BaseHandler = util.NewBaseHandler()

	ws, err := util.NewWebsocket(w, r)
	if err != nil {
		base.Error(err.Error())
		util.GlobalEvents.Shutdown(1)
	}

	h := Handler{BaseHandler: base, WebSocket: ws}

	ws.HandleConnection(h.listenToWebSocket)

}
