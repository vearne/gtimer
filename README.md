# timer
用golang实现的定时器，基于delayqueue

# 实现
实现受到了Java DelayQueue.java的启发
源码地址
[DelayQueue.java](http://www.docjar.com/html/api/java/util/concurrent/DelayQueue.java.html)

依赖的几个结构依次为为
timer -> delayqueue -> priorityqueue -> heap

由于golang的Condition不支持wait一段时间，所以使用golang原生的Timer来替代了Condition在delayqueue中的作用

# Installation
## Install:

```
go get -u github.com/vearne/gtimer
```
## Import:
```
import "github.com/vearne/gtimer"
```


## Quick Start
```
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vearne/gtimer"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	PRODUCER_COUNT = 10
	CONSUMER_COUNT = 10
	TARGET_COUNT   = 1000000
)

var ops int64 = 0

func main() {
	st := gtimer.NewSuperTimer(CONSUMER_COUNT, time.Second)
	// st := gtimer.NewSuperTimer(CONSUMER_COUNT, 0)
	t1 := time.Now()
	for i := 0; i < PRODUCER_COUNT; i++ {
		go push(st, "worker"+strconv.Itoa(i))
	}

	time.Sleep(100 * time.Millisecond)

	for {
		v := atomic.LoadInt64(&ops)
		if v >= TARGET_COUNT {
			st.Stop()
			break
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
	t2 := time.Now()
	log.Infof("cost:%v\n", t2.Sub(t1))
}

func DefaultAction(t time.Time, value string) {
	// fmt.Printf("trigger_time:%v, value:%v\n", t, value)
	atomic.AddInt64(&ops, 1)
}

func push(timer *gtimer.SuperTimer, name string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000000; i++ {
		now := time.Now()
		t := now.Add(time.Millisecond * time.Duration(r.Int63n(300)))
		value := fmt.Sprintf("%v:value:%v", name, strconv.Itoa(i))
		// create a delayed task
		item := gtimer.NewDelayedItemFunc(t, value, DefaultAction)
		timer.Add(item)
	}
}
```

use NewDelayedItemFunc, we can create a task
```
// triggerTime is time of the task should be execute
func NewDelayedItemFunc(triggerTime time.Time, value string, f func(time.Time, string)) *Item
```
task struct like 
```
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int64    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
	// when task is ready, execute OnTrigger function
	OnTrigger func(time.Time, string)
}
```

## Performance

`CPU Model Name`: 2.3 GHz Intel Core i5     
`CPU Processors`: 4     
`Memory`: 8GB   

### Benchmark Test Results        


| produce goroutines count | consume goroutines count | qps(per second) | 
| ---------:| ----------:| --------:| 
| 1| 1                      | 285714             |  
| 10| 10                    | 90090                |  
| 10| 100                   | 89285              |  
| 100| 100                  | 23255              |  





