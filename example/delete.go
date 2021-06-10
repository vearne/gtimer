package main

import (
	"fmt"
	"github.com/vearne/gtimer"
	"time"
)

func main() {
	startTime := time.Now()
	gtimer.Add(3*time.Second, func() {
		fmt.Println("task1", time.Since(startTime))
	})
	task2 := gtimer.Add(5*time.Second, func() {
		fmt.Println("task2", time.Since(startTime))
	})
	gtimer.Add(7*time.Second, func() {
		fmt.Println("task3", time.Since(startTime))
	})
	fmt.Println("before remove, size:", gtimer.Size())
	go func() {
		time.Sleep(10 * time.Second)
		gtimer.Stop()
	}()

	gtimer.Remove(task2)
	fmt.Println("after remove, size:", gtimer.Size())
	gtimer.Wait()
	fmt.Println("---end---")
}
