package util

import (
	"fmt"
	"sync"
)

type ReadWriteLock struct {
	ReadLock  sync.Mutex
	WriteLock sync.Mutex
}

type Uid string

func NewUid() Uid {
	return Uid(GenerateUid())
}

func (u Uid) GetUid() string {
	return string(u)
}

type BaseHandler struct {
	Uid
	Logger
}

func (h *BaseHandler) Initialize() {
	if h.GetUid() == "" {
		h.Uid = NewUid()
		h.Logger = NewLogger(fmt.Sprintf("[%s]", h.Uid[:8]))
	}

}
