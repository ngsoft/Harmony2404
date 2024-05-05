package wsocket

type ChannelHandler struct {
	Channels map[EventType]*Channel
	Clients  map[*Client]*Channel
	init     bool
}

func (h *ChannelHandler) __init() {
	if !h.init {
		h.init = true
		h.Channels = make(map[EventType]*Channel)
		h.Clients = make(map[*Client]*Channel)
	}
}

func (h *ChannelHandler) AddChannel(ch *Channel, ev EventType) {
	h.__init()
	h.Channels[ev] = ch
}

func (h *ChannelHandler) OnMessageEvent(c *Client, e *Event) {
	h.__init()

	if e.Direction == Incoming {
		t := e.Type
		if !e.Channel {
			// channel bound to the event
			if ch, ok := h.Channels[t]; ok {
				// switching channel
				currentChan, ok := h.Clients[c]
				if ok && currentChan != ch {
					currentChan.UnRegisterClient(c)
					ch.RegisterClient(c)
				}
				if !ok {
					ch.RegisterClient(c)
				}
				// keep a record of the client channel
				h.Clients[c] = ch
				e.Channel = true
				// add event to channel
				ch.BroadcastEvent(e)
				return
			}
			// client is in a channel
			if ch, ok := h.Clients[c]; ok {
				// let the channel handle the event
				ch.BroadcastEvent(e)
				return
			}

			// client not in channel
			ev := NewEvent(ErrorEvent, "invalid event sent")
			c.SendEvent(&ev)
		}

	}

}
func (h *ChannelHandler) OnMessage(c *Client, m string) {
	if ch, ok := h.Clients[c]; ok {
		ch.BroadcastMessage(m)
	}
}
func (h *ChannelHandler) OnOpen(c *Client) {
	h.__init()

}
func (*ChannelHandler) OnClose(c *Client) {
	// noop
}
