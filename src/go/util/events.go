package util

type EventType string

func (t EventType) Is(tt EventType) bool {
	return tt == t
}

func (t EventType) String() string {
	return string(t)
}

const (
	allEvents       EventType = "util.all"
	SignalEvent     EventType = "util.signal"
	InitializeEvent EventType = "util.init"
	ShutdownEvent   EventType = "util.shutdown"
)

type Event struct {
	Type   EventType
	Params []interface{}
	Status
}

func (h *Event) StopPropagation() {
	h.Status = STATUS_OFF
}

type EventHandler interface {
	OnEvent(*Event)
}

type EventListener struct {
	init       bool
	handlers   map[EventType]map[EventHandler]bool
	registered map[EventHandler]int
}

func (h *EventListener) __init() {
	if !h.init {
		h.init = true
		h.handlers = make(map[EventType]map[EventHandler]bool)
		h.registered = make(map[EventHandler]int)
	}
}

func (h *EventListener) AddEventHandler(e EventHandler, v ...EventType) func() {
	h.__init()

	var (
		id int
		ok bool
	)

	if _, ok = h.registered[e]; !ok {
		id = len(h.registered)
		h.registered[e] = id
	}

	if len(v) == 0 {
		v = append(v, allEvents)
	}

	for _, t := range v {
		if _, ok := h.handlers[t]; !ok {
			h.handlers[t] = make(map[EventHandler]bool)
		}
		h.handlers[t][e] = true
	}

	return func() {
		h.RemoveEventHandler(e, v...)
	}
}

func (h *EventListener) RemoveEventHandler(e EventHandler, v ...EventType) {
	if !h.init {
		return
	}

	if len(v) == 0 {
		v = append(v, allEvents)
	}

	for _, t := range v {
		if _, ok := h.handlers[t]; ok {
			delete(h.handlers[t], e)
		}
	}

}

func (h *EventListener) DispatchEvent(e EventType, p ...interface{}) {
	if !h.init {
		return
	}
	var (
		ev = Event{
			Type:   e,
			Params: p,
			Status: STATUS_ON,
		}
	)

	for l := range h.registered {
		if !ev.On() {
			break
		}
		if _, ok := h.handlers[allEvents]; ok {
			l.OnEvent(&ev)
			continue
		}
		if _, ok := h.handlers[e]; ok {
			l.OnEvent(&ev)
		}
	}
}
