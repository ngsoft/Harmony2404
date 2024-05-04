package main

import (
	"flirc/util"
	"fmt"
)

type StatusType string

const (
	StatusSuccess StatusType = "ok"
	StatusError   StatusType = "ko"
)

type BaseResponse struct {
	Status  StatusType  `json:"status"`
	Message string      `json:"statusMessage"`
	Data    interface{} `json:"result"`
}

type BaseEvent struct {
	Type   string      `json:"type"`
	Params interface{} `json:"params"`
}

type SocketMessage struct {
	Message string
	Code    string
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
	Message SocketMessage
	remote  string
	path    string
	util.EventListener
	util.BaseHandler
}
