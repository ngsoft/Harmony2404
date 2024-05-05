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
		// channel exists
		if newChan, ok := h.Channels[t]; ok {
			// switching channel
			if currentChan, ok := h.Clients[c]; ok {
				if currentChan != newChan {
					currentChan.UnRegisterClient(c)
					newChan.RegisterClient(c)
				}
			}
			// add event to channel
			newChan.BroadcastEvent(e)
			return
		}
		// not existing channel
		ev := NewEvent(ErrorEvent, "invalid event sent")
		c.SendEvent(&ev)
	}

}
func (*ChannelHandler) OnMessage(c *Client, m string) {
	ev := NewEvent(ErrorEvent, "invalid message sent")
	c.SendEvent(&ev)
}
func (h *ChannelHandler) OnOpen(c *Client) {
	h.__init()

}
func (*ChannelHandler) OnClose(c *Client) {

}
