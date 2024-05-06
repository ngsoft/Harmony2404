package util

var (
	globalEvents EventListener
)

type EventType string

func (t EventType) Is(tt EventType) bool {
	return tt == t
}

func (e EventType) String() string {
	return string(e)
}

const (
	allEvents       EventType = "all"
	SignalEvent     EventType = "signal"
	InitializeEvent EventType = "init"
	ShutdownEvent   EventType = "shutdown"
)

type EventDisabler func()

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

func (h *EventListener) AddEventHandler(e EventHandler, v ...EventType) EventDisabler {
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
func AddEventHandler(e EventHandler, v ...EventType) EventDisabler {
	return globalEvents.AddEventHandler(e, v...)
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
func RemoveEventHandler(e EventHandler, v ...EventType) {
	globalEvents.RemoveEventHandler(e, v...)
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
func DispatchEvent(e EventType, p ...interface{}) {
	globalEvents.DispatchEvent(e, p...)
}
