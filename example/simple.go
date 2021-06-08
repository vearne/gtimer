package main

import (
	"fmt"
	"github.com/vearne/gtimer"
	"time"
)

func main() {
	startTime := time.Now()
	gtimer.Add(3*time.Second, func() {
		fmt.Println(time.Since(startTime))
	})
	go func() {
		time.Sleep(5 * time.Second)
		gtimer.Stop()
	}()
	gtimer.Wait()
	fmt.Println("---end---")
}
