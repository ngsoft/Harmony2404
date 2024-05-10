package main

import (
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
