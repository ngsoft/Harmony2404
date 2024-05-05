package wsocket

import (
	"encoding/json"
)

type Direction int

const (
	Incoming Direction = iota
	Outgoing
)

type EventType string

type Payload []interface{}

type Event struct {
	Direction
	Type EventType
	Data interface{}
}

func (m Event) Bytes() []byte {
	payload := m.Payload()
	b, err := json.Marshal(payload)
	if err != nil {
		return make([]byte, 0)
	}
	return b
}
func (m Event) String() string {
	return string(m.Bytes())
}

func (m Event) Payload() Payload {
	return Payload{
		m.Type,
		m.Data,
	}
}
func NewEventFromString(s string) (Event, bool) {
	var (
		p Payload
		e Event
	)
	if err := json.Unmarshal([]byte(s), &p); err == nil {
		return NewEventFromPayload(p)
	}
	return e, false
}
func NewEventFromPayload(p Payload) (Event, bool) {
	var e Event
	if len(p) > 0 {
		if v, ok := p[0].(string); ok {
			e.Type = EventType(v)
			e.Data = nil
			if len(p) > 1 {
				e.Data = p[1]
			}
			return e, true
		}
	}
	return e, false
}

func NewEvent(t EventType, d ...interface{}) Event {
	var data interface{} = nil
	if len(d) > 0 {
		data = d[0]
	}
	return Event{Type: t, Data: data}
}
