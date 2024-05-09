package main

import (
	"flirc/usocket"
	"flirc/util"
	"flirc/wsocket"
	"fmt"
)

type Keymap string
type State string

const (
	StateSuccess State        = "ok"
	StateError   State        = "ko"
	Input        wsocket.Type = "input"
)

func (s State) OK() bool {
	return s == StateSuccess
}

type JsonResponse struct {
	State   `json:"status"`
	Message string `json:"statusMessage"`
}

type JsonResponseWithData struct {
	*JsonResponse
	Data interface{} `json:"result"`
}

type BaseEvent struct {
	Type   string      `json:"type"`
	Params interface{} `json:"params"`
}

type SocketMessage struct {
	Message string
	Code    string
	CodeInt int
	Repeat  int
	Key     string
	Remote  string
}

func (h *SocketMessage) Export() []interface{} {
	return []interface{}{
		fmt.Sprintf("0x%s", h.Code),
		h.Key,
		h.Repeat,
	}
}

type Handler struct {
	util.BaseHandler
}

type FlircHandler struct {
	keymaps []Keymap
	remote  string
	delay   int
	*usocket.UnixSocket
	util.Logger
	lock bool
	Room *wsocket.Room
}
