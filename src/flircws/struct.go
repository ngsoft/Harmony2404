package main

import (
	"flirc/usocket"
	"flirc/util"
	"fmt"
)

type Keymap string

type Status string

const (
	StatusSuccess Status = "ok"
	StatusError   Status = "ko"
)

func (s Status) OK() bool {
	return s == StatusSuccess
}

type JsonResponse struct {
	Status  `json:"status"`
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

type ConnHandler struct {
	util.Logger
	handlers map[string]Handler
	Route    string
	Port     int
}

func (c *ConnHandler) GetHttpPort() string {
	return fmt.Sprintf(":%d", c.Port)
}

type Handler struct {
	util.BaseHandler
}

const (
	INPUT_EVENT util.EventType = "flirc_input"
)

type FlircHandler struct {
	keymaps []Keymap
	remote  string
	delay   int
	*usocket.UnixSocket
	util.Logger
	lock bool
}
