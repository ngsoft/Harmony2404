package wsocket

import (
	"flirc/util"
	"time"
)

const (
	Login Type = "login"
)

type Username interface{}
type Password interface{}

type User struct {
	name     Username
	password Password
}

func NewUser(name Username, password Password) *User {
	return &User{name, password}
}

type UserValidator interface {
	ValidateUser(Username, Password) bool
}

type InMemoryUserValidator struct {
	users map[Username]*User
	init  bool
}

func (v *InMemoryUserValidator) __init() {
	if !v.init {
		v.init = true
		v.users = make(map[Username]*User)
	}
}

func (l *InMemoryUserValidator) AddUsers(u *User, v ...*User) {
	l.__init()
	l.users[u.name] = u
	for _, user := range v {
		l.users[user.name] = user
	}
}

func (v *InMemoryUserValidator) ValidateUser(u Username, p Password) bool {
	v.__init()
	if user, ok := v.users[u]; ok {
		return p == user.password
	}
	return false
}

type LoginHandler struct {
	validator UserValidator
	clients   map[*Client]bool
	timeouts  map[*Client]func()
	timeout   time.Duration
	init      bool
	ws        *WebSocket
}

func NewLoginHandler(ws *WebSocket, timeout uint, validator UserValidator) *LoginHandler {

	if timeout == 0 {
		panic("timeout cannot be 0 seconds")
	}

	var l = LoginHandler{
		validator: validator,
		clients:   make(map[*Client]bool),
		timeouts:  make(map[*Client]func()),
		timeout:   time.Duration(timeout) * time.Second,
		init:      true,
		ws:        ws,
	}
	ws.AddEventHandler(&l)
	return &l
}

func (l *LoginHandler) __init() {
	if !l.init {
		panic("Please use NewLoginHandler()")
	}
}

func (l *LoginHandler) OnEvent(ev *util.Event) {
	l.__init()

	switch ev.Type {
	case Open, Close:
		if len(ev.Params) > 0 {
			if client, ok := ev.Params[0].(*Client); ok {
				if ev.Type == Close {
					if _, ok := l.timeouts[client]; ok {
						l.timeouts[client]()
						delete(l.timeouts, client)
					}
					delete(l.clients, client)
					return
				}
				l.timeouts[client] = util.SetTimeout(func() {
					client.Error("login timeout reached (%ds)", int(l.timeout/time.Second))
					client.SendEvent(Error, Login)
					client.Close()
				}, l.timeout)
			}
		}
	}

}

func (l *LoginHandler) OnMessage(m *MessageEvent, next *NextHandler) {
	l.__init()
	var (
		connected bool
		client    = m.Client
	)
	if _, connected = l.clients[client]; connected {
		next.OnMessage(m)
	}
	if m.Direction.IsIncoming() && !connected {
		if m.Type == Login {
			if len(m.Params) >= 2 {
				if name, ok := m.Params[0].(string); ok {
					name := Username(name)
					if pass, ok := m.Params[1].(string); ok {
						pass := Password(pass)
						if l.validator.ValidateUser(name, pass) {
							l.timeouts[client]()
							l.clients[client] = true
							client.Info("logged in as %s", name)
							client.SendEvent(Success, Login, name)
							return

						}
					}
				}
			}
			client.Error("invalid credentials provided")
			client.SendEvent(Error, Login)
			client.Close()
		}
	}
}

func (l *LoginHandler) SetUserValidator(u UserValidator) {
	l.validator = u
}
