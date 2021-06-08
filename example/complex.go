package main

import (
	"fmt"
	"github.com/vearne/gtimer"
	"time"
)

func doSomeThing(scheduledExecTime time.Time, param interface{}) {
	fmt.Println("Scheduled Execution Time:", scheduledExecTime, "Real Execution Time:", time.Now())
	fmt.Println("hello ", param)
}

func main() {
	startTime := time.Now()
	// change default timer
	//gtimer.Set(gtimer.NewSuperTimer(20))
	gtimer.AddComplex(startTime.Add(3*time.Second), doSomeThing, "world")
	gtimer.AddComplex(startTime.Add(4*time.Second), doSomeThing, "sky")
	gtimer.AddComplex(startTime.Add(2*time.Second), doSomeThing, "groud")
	go func() {
		time.Sleep(5 * time.Second)
		gtimer.Stop()
	}()
	gtimer.Wait()
	fmt.Println("---end---")
}
