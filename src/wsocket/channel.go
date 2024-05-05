package wsocket

import "flirc/util"

type ChannelMessageHandler interface {
	OnChannelEvent(*Channel, *Event)
	OnChannelMessage(*Channel, string)
}

type Channel struct {
	util.BaseHandler
	Handler ChannelMessageHandler
	Clients map[*Client]bool
	event   chan *Event
	message chan string
	add     chan *Client
	remove  chan *Client
}

func CreateChannel() *Channel {
	c := Channel{
		event:   make(chan *Event),
		message: make(chan string),
		add:     make(chan *Client),
		remove:  make(chan *Client),
		Clients: make(map[*Client]bool),
	}
	c.Handler = &c
	c.Initialize()
	return &c
}

func (h *Channel) Run() {
	for {
		select {
		case client := <-h.add:
			h.Clients[client] = true
		case client := <-h.remove:
			delete(h.Clients, client)
		case event := <-h.event:
			h.Handler.OnChannelEvent(h, event)
		case msg := <-h.message:
			h.Handler.OnChannelMessage(h, msg)
		}

	}
}

func (h *Channel) OnChannelEvent(c *Channel, event *Event) {
	for client := range h.Clients {
		select {
		case client.incoming <- event:
		default:
			delete(h.Clients, client)
		}
	}
}
func (h *Channel) OnChannelMessage(c *Channel, m string) {
	// noop
}

func (h *Channel) SetHandler(v ChannelMessageHandler) {
	h.Handler = v
}
func (h *Channel) RegisterClient(c *Client) {
	if _, ok := h.Clients[c]; !ok {
		h.add <- c
	}
}

func (h *Channel) UnRegisterClient(c *Client) {
	h.remove <- c
}
func (h *Channel) BroadcastEvent(e *Event) {
	h.event <- e
}
