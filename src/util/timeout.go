package util

import (
	"time"

	"github.com/imbhargav5/noop"
)

func SetTimeout(task func(), duration time.Duration) func() {

	var active = true
	timer := *time.AfterFunc(duration, func() {
		if !active {
			return
		}
		active = false
		task()
	})

	return func() {
		task = noop.Noop
		active = false
		if !timer.Stop() {
			<-timer.C
		}
	}
}
