package wsocket

const (
	JoinRoom Type = "join"
)

type Room struct {
	Name      string
	Clients   map[*Client]bool
	WebSocket *WebSocket
	Handler   Handler
}

func (r *Room) AddClient(c *Client) {

	if c.On() {
		if c.CurrentRoom != nil {
			c.CurrentRoom.RemoveClient(c)
		}
		r.Clients[c] = true
		c.CurrentRoom = r
		c.SendEvent(Success, JoinRoom, r.Name)
	}

}

func (r *Room) RemoveClient(c *Client) {
	c.CurrentRoom = nil
	delete(r.Clients, c)
}

func (r *Room) SetHandler(h Handler) {
	r.Handler = h
}
func (r *Room) OnMessage(c *Client, d Direction, t Type, v ...interface{}) {

	if r.Handler != nil {
		r.Handler.OnMessage(c, d, t, v...)
		return
	}

	// basic room broadcast for everyone except sender
	if d.IsIncoming() {
		for cl := range r.Clients {
			if cl != c {
				cl.SendEvent(t, v...)
			}
		}
	}

}

// SendEvent send event to everyone in the room
func (r *Room) SendEvent(t Type, v ...interface{}) {
	for c := range r.Clients {
		c.SendEvent(t, v...)
	}
}
