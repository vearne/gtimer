package gtimer

import "time"

var (
	defaultTimer *SuperTimer
)

func init() {
	defaultTimer = NewSuperTimer(10)
}

func Set(timer *SuperTimer) {
	defaultTimer = timer
}

func Add(d time.Duration, f func()) {
	task := NewDelayedItemFunc(time.Now().Add(d), nil, func(t time.Time, i interface{}) {
		f()
	})
	defaultTimer.Add(task)
}

func AddComplex(scheduledExecTime time.Time, f func(scheduledExecTime time.Time, param interface{}), param interface{}) {
	task := NewDelayedItemFunc(scheduledExecTime, param, f)
	defaultTimer.Add(task)
}

func Stop() {
	defaultTimer.Stop()
}

func Wait() {
	defaultTimer.Wait()
}
