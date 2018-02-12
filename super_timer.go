package gtimer

import (
	"time"
	"sync"
	"container/heap"
	log "github.com/sirupsen/logrus"
)

type SuperTimer struct{
	lock *sync.Mutex
	PQ PriorityQueue
	UniTimer *time.Timer
	WorkerCount int
	Wgp *sync.WaitGroup
	ExitChan chan int
	RunningFlag bool
}


func NewSuperTimer(workCount int)*SuperTimer{
	timer := SuperTimer{}
	timer.lock = &sync.Mutex{}
	timer.UniTimer = time.NewTimer(time.Second * 10)
	timer.PQ = make(PriorityQueue, 0, 100)
	timer.WorkerCount = workCount
	timer.Wgp = &sync.WaitGroup{}
	timer.ExitChan = make(chan int, 100)
	timer.RunningFlag = true

	for i :=0; i < timer.WorkerCount; i++ {
		go func(){
			timer.Wgp.Add(1)
			defer timer.Wgp.Done()
			for ;timer.RunningFlag; {
				select {
					case <- timer.ExitChan:
						log.Debug("ExitChan")
						break
					case <- timer.UniTimer.C:
						item := timer.Take()
						log.Debugf("[consumer]%v\n", item)
						if item != nil {
							log.Debugf("[consumer]%v, item:%v, %v\n", time.Now(), item.priority, item.value)
							t := time.Unix(item.priority/1000000000, item.priority%1000000000)
							item.OnTrigger(t, item.value)
						}
				}
			}
			log.Debugf("worker exit")
		}()
	}
	return &timer
}

func (st *SuperTimer) Add(pItem *Item){
	st.lock.Lock()
	defer st.lock.Unlock()
	heap.Push(&st.PQ, pItem)
	log.Infof("[producer] PQ size:%v", len(st.PQ))
	st.PQ.Dump()
	peek := st.PQ.Peek().(*Item)
	if peek == pItem{
		log.Infof("[producer] reset:%v", pItem.GetDelay())
		st.UniTimer.Reset(pItem.GetDelay())
	}
}

func (st *SuperTimer) Take() *Item{
	st.lock.Lock()
	defer st.lock.Unlock()
	log.Infof("[producer] PQ size:%v", len(st.PQ))
	st.PQ.Dump()
	if len(st.PQ) <= 0 {
		st.UniTimer.Reset(time.Second * 5)
		log.Infof("[consumer]reset %v", 5)
		return nil
	}

	item := st.PQ.Peek()
	target:=item.(*Item);
	if target.GetDelay() > 0{
		st.UniTimer.Reset(target.GetDelay())
		log.Infof("[consumer]reset %v, target:%v, now:%v", target.GetDelay(), target.priority, time.Now().UnixNano())
		return nil
	}
	
	res := heap.Pop(&st.PQ).(*Item)
	// 重置定时器，立刻唤醒其它消费者
	st.UniTimer.Reset(0)
	return res
}

func (st *SuperTimer) Stop(){
	st.lock.Lock()
	defer st.lock.Unlock()
	log.Infof("------- stop ------")
	// 强制工作线程退出
	st.PQ.Clear()
	st.RunningFlag = false
	close(st.ExitChan)
	st.UniTimer.Stop()
}

func (st *SuperTimer) Wait(){
	st.Wgp.Wait()
}

