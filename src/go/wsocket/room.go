package wsocket

const (
	JoinRoom  Type = "join"
	LeaveRoom Type = "leave"
)

type Room struct {
	Name      string
	Clients   map[*Client]bool
	WebSocket *WebSocket
}

func (r *Room) AddClient(c *Client) {

	if c.On() {
		if c.CurrentRoom != nil {
			c.CurrentRoom.RemoveClient(c)
		}
		r.Clients[c] = true
		c.CurrentRoom = r
	}

}

func (r *Room) RemoveClient(c *Client) {
	c.CurrentRoom = nil
	delete(r.Clients, c)
}

// SendEvent send event to everyone in the room
func (r *Room) SendEvent(t Type, v ...interface{}) {
	for c := range r.Clients {
		c.SendEvent(t, v...)
	}
}
