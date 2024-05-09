package util

import (
	"time"
)

func SetTimeout(task func(), duration time.Duration) func() {

	var (
		cnt   int
		timer = time.NewTimer(duration)
	)

	go func() {
		<-timer.C
		if cnt == 0 {
			task()
		}
	}()

	return func() {
		cnt++
		if !timer.Stop() {
			<-timer.C
		}
	}
}

func SetInterval(task func(), duration time.Duration) func() {

	var (
		active = true
		ticker = time.NewTicker(duration)
	)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if !active {
				return
			}
			task()
		}
	}()

	return func() {
		active = false
		ticker.Stop()
	}

}
