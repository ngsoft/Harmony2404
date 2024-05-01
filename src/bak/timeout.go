package main

import (
	"time"

	"github.com/imbhargav5/noop"
)

type Timeout struct {
	uid      string
	callback func()
	timer    time.Timer
	active   bool
}

func (t *Timeout) Stop() {
	t.active = false
	t.callback = noop.Noop
	if !t.timer.Stop() {
		<-t.timer.C
	}
}

func setTimeout(callback func(), duration int) Timeout {

	t := Timeout{
		uid:      generateUid(),
		callback: callback,
		active:   true,
	}

	t.timer = *time.AfterFunc(time.Duration(int64(duration))*time.Millisecond, func() {
		if !t.active {
			return
		}
		t.callback()
		t.active = false
	})

	return t
}
