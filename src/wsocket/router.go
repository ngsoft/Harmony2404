package wsocket

import (
	"flirc/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

const (
	DefaultPingInterval = 60 * time.Second
	DefaultPingEnabled  = false
)

type RouterConfig struct {
	Route         string
	ClientHandler ClientHandler
	PingEnabled   bool
	PingInterval  time.Duration
}

type Router struct {
	*RouterConfig
	util.Logger
	Clients map[*Client]bool
}

func (r *Router) OnEvent(ev *util.Event) {

	switch ev.Type {
	case util.ShutdownEvent:
		// close all clients connections
		r.Info("received shutdown event, closing all clients connections")
		for client := range r.Clients {
			if client.On() {
				client.Close()
			}
		}

	}

}

func NewRouterFromConfig(c RouterConfig) Router {
	return Router{
		Logger:       util.NewLogger(fmt.Sprintf("[route:%s]", c.Route)),
		RouterConfig: &c,
		Clients:      make(map[*Client]bool),
	}
}

func NewRouter(route string, v ...interface{}) Router {
	var (
		cfg RouterConfig
		p                 = DefaultPingEnabled
		i                 = DefaultPingInterval
		h   ClientHandler = nil
	)

	for _, param := range v {
		if val, ok := param.(bool); ok {
			p = val
			continue
		}
		if val, ok := param.(time.Duration); ok {
			i = val
			continue
		}

		if val, ok := param.(ClientHandler); ok {
			h = val
		}
	}
	cfg.PingEnabled = p
	cfg.PingInterval = i
	cfg.ClientHandler = h
	return NewRouterFromConfig(cfg)
}

func (rt *Router) SetClientHandler(h ClientHandler) {
	rt.RouterConfig.ClientHandler = h
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		rt.Error(err.Error())
		return
	}

	rt.Info("new connection, executing handler")
	client := NewClient(rt, conn)
	rt.Clients[client] = true
	util.AddEventHandler(rt, util.ShutdownEvent)

	client.Run()

}
