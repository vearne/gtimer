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
	"time"
	"github.com/vearne/gtimer"
	"strconv"
	"math/rand"
	"sync"
	log "github.com/sirupsen/logrus"
)



func main(){
	t1 := time.Now()
	wg := sync.WaitGroup{}
	timer := gtimer.NewSuperTimer(1)
	// concurrent push task
	for i:=0;i<1;i++{
		wg.Add(1)
		go push(timer, "worker" + strconv.Itoa(i))
		wg.Done()
	}
	// wg.Wait()
	log.Infof("[producer]------push ok-------")
	go func(){
		log.Infof("[start]try to stop")
		time.Sleep(5 * time.Second)
		for {
			if timer.Size() > 0{
				time.Sleep(1)
			}else{
				break
			}
		}
		timer.Stop()
		log.Infof("[end]try to stop")
	}()
	// wait until stop
	timer.Wait()
	t2 := time.Now()
	log.Infof("cost:%v\n", t2.Sub(t1))

}

func DefaultAction(t time.Time, value string){
	fmt.Printf("trigger_time:%v, value:%v\n", t, value)
}

func push(timer *gtimer.SuperTimer, name string){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0;i< 10000;i++{
		now := time.Now()
		t := now.Add(time.Millisecond * time.Duration(r.Int63n(300)) + 1 * time.Second)
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
CPU Model Name: 2.3 GHz Intel Core i5
CPU Processors: 4
Memory: 8GB

### Test Results:
|produce goroutines count|consume goroutines count|qps(per second)|
|:---|:---|:---|:---|
|1|1|10000|
|5|1|10000|
|1|5|10000|
|5|5|10000|


