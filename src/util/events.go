package util

type EventType string

func (t EventType) Is(tt EventType) bool {
	return tt == t
}

const (
	allEvents      EventType = "all"
	SIGNAL_EVENT   EventType = "signal"
	SHUTDOWN_EVENT EventType = "shutdown"
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
	HandleEvent(*Event)
}

type EventListener struct {
	init     bool
	handlers map[string]map[EventType]EventHandler
}

func (h *EventListener) __init() {
	if !h.init {
		h.init = true
		h.handlers = make(map[string]map[EventType]EventHandler)
	}
}

func (h *EventListener) AddEventHandler(e EventHandler, v ...EventType) EventDisabler {
	h.__init()
	uid := GenerateUid()

	if len(v) == 0 {
		v = append(v, allEvents)
	}

	for _, t := range v {

		if _, ok := h.handlers[uid]; !ok {
			h.handlers[uid] = make(map[EventType]EventHandler)
		}
		h.handlers[uid][t] = e
	}

	return func() {
		delete(h.handlers, uid)
	}
}

func (h *EventListener) DispatchEvent(e EventType, p ...interface{}) bool {
	if !h.init || len(h.handlers) == 0 {
		return false
	}
	var (
		result bool
		ev     = Event{
			Type:   e,
			Params: p,
			Status: STATUS_ON,
		}
	)

	for _, m := range h.handlers {
		if !ev.On() {
			break
		}
		if c, ok := m[allEvents]; ok {
			result = true
			c.HandleEvent(&ev)
		} else if c, ok := m[e]; ok {
			result = true
			c.HandleEvent(&ev)
		}
	}

	return result
}
