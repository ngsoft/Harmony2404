package util

import (
	"fmt"
	"sync"
)

type ReadWriteLock struct {
	ReadLock  sync.Mutex
	WriteLock sync.Mutex
}

type BaseHandler struct {
	Uid string
	Logger
}

func (h *BaseHandler) Initialize() {
	if h.Uid == "" {
		h.Uid = GenerateUid()
		h.Logger = NewLogger(fmt.Sprintf("[%s]", h.Uid[:10]))
	}

}
