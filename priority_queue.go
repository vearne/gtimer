// This example demonstrates a priority queue built using the heap interface.
package gtimer

import (
	"container/heap"
	"time"
	log "github.com/sirupsen/logrus"
)


// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n - 1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int64) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) Peek() interface{} {
	old := *pq
	n := len(old)
	if n <= 0{
		return nil
	}
	return old[0]
}

func (pq *PriorityQueue) Clear(){
	pq = nil
}

func (pq *PriorityQueue) Dump(){
	old := *pq
	log.Debugf("---------------------")
	for i, v:= range old{
		log.Debugf("index:%v, item:%v, delayed:%v\n", i, v, (v.priority - time.Now().UnixNano()) / 1000000000)
	}
	log.Debugf("---------------------")
}
