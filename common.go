package gtimer

import (
	"time"
)


type Delayed interface{
	GetDelay() time.Duration
}


// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int64    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
	OnTrigger func(time.Time, string)
}

func NewDelayedItemFunc(triggerTime time.Time, value string, f func(time.Time, string)) *Item{
	item := Item{}
	item.priority = triggerTime.UnixNano()
	item.value = value
	item.OnTrigger = f
	return &item
}

func (item Item) GetDelay() time.Duration{
	return time.Duration(item.priority - time.Now().UnixNano())
}






