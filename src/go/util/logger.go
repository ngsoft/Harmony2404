package util

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger struct {
	*log.Logger
	prefix string
	suffix string
}

func NewLogger(v ...interface{}) Logger {

	var l Logger
	for _, arg := range v {
		if n, ok := arg.(Logger); ok {
			l = n
			continue
		}
		if n, ok := arg.(string); ok {
			if l.prefix == "" {
				l.prefix = strings.TrimRight(n, " ") + " "
				continue
			}
			l.suffix += " " + strings.TrimLeft(n, " ")
		}
	}
	return l
}
func (h *Logger) SetLoggerPrefix(pfx string) {
	h.prefix = strings.TrimRight(pfx, " ") + " "
}
func (h *Logger) __init() {
	if h.Logger == nil {
		h.Logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	}
}
func (h *Logger) Log(m string, v ...interface{}) {
	h.__init()
	_ = h.Logger.Output(2, fmt.Sprintf(h.prefix+m+h.suffix, v...))
}
func (h *Logger) Warn(m string, v ...interface{}) {
	h.__init()
	_ = h.Logger.Output(2, fmt.Sprintf(h.prefix+"WARN: "+m+h.suffix, v...))
}
func (h *Logger) Error(m string, v ...interface{}) {
	h.__init()
	_ = h.Logger.Output(2, fmt.Sprintf(h.prefix+"ERR: "+m+h.suffix, v...))
}
func (h *Logger) Info(m string, v ...interface{}) {
	h.__init()
	_ = h.Logger.Output(2, fmt.Sprintf(h.prefix+"INFO: "+m+h.suffix, v...))
}
