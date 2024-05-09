package wsocket

import (
	"flirc/util"
)

const (
	Open  util.EventType = "wsocket.open"
	Close util.EventType = "wsocket.close"
)

type Type string

const (
	Success Type = "success"
	Error   Type = "error"
	Message Type = "message"
)

type Param interface{}
type Payload []interface{}

func (e Type) String() string {
	return string(e)
}

func (e Type) Is(v interface{}) bool {
	if val, ok := v.(Type); ok {
		return val == e
	}
	return false
}

type Direction bool

const (
	In  Direction = true
	Out Direction = false
)

func (d Direction) IsIncoming() bool {
	return d == In
}

func (d Direction) String() string {
	if d.IsIncoming() {
		return "in"
	}
	return "out"
}

type Handler interface {
	OnMessage(*Client, Direction, Type, ...interface{})
}
