package main

import (
	"flirc/util"
	"flirc/wsocket"
	"strings"
)

// triggers
const (
	GetInput   wsocket.Type = "getInputEvents"
	GetKeymaps wsocket.Type = "getKeymaps"
	Input      wsocket.Type = "input"
)

type KeyPair struct {
	Keycode `json:"code"`
	Keymap  `json:"key"`
}

func (k KeyPair) List() []interface{} {
	return []interface{}{
		k.Keycode, k.Keymap,
	}
}

type InputEvent struct {
	KeyPair
	Remote string `json:"-"`
}

func NewInputEvent(m string) (InputEvent, bool) {

	if len(m) > 0 {
		parts := strings.Split(m, " ")

		if len(parts) == 4 {
			return InputEvent{
				Remote: parts[3],
				KeyPair: KeyPair{
					Keycode: Keycode(util.HexToInt(parts[0])),
					Keymap:  Keymap(parts[2]),
				},
			}, true
		}

	}
	return InputEvent{}, false
}
