package util

import (
	"net/http"
)

type RouteHandler struct {
	route          string
	RequestHandler func(http.ResponseWriter, *http.Request)
}

func NewRouteHandler(route string) (RouteHandler, bool) {
	var rt RouteHandler
	if route[:1] != "/" {
		return rt, false
	}
	rt = RouteHandler{route: route}
	http.Handle(route, &rt)
	return rt, true
}

func (h *RouteHandler) IsValid() bool {
	return h.route[:1] == "/"
}

func (h *RouteHandler) SetRequestHandler(fn func(http.ResponseWriter, *http.Request)) bool {

	if !h.IsValid() {
		return false
	}
	h.RequestHandler = fn
	return true
}

func (h *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if isset := h.RequestHandler; isset != nil {
		h.RequestHandler(w, r)
	}
}
