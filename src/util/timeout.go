package util

import (
	"time"

	"github.com/imbhargav5/noop"
)

type TimeoutEnder func()

func SetTimeout(task func(), duration int) TimeoutEnder {

	var active bool = true
	timer := *time.AfterFunc(time.Duration(int64(duration))*time.Millisecond, func() {
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
