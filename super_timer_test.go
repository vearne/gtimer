package gtimer


import (
	"testing"
	"fmt"
	"time"
	"strconv"
	"math/rand"
)

func push(timer *SuperTimer, name string){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0;i<20;i++{
		now := time.Now()
		t := now.Add(time.Duration(r.Int63n(100)) * time.Millisecond)
		value := fmt.Sprintf("%v:value:%v", name, strconv.Itoa(i))
		item := NewDelayedItem(t, value)
		timer.Add(item)
	}
}

func Test_timer(t *testing.T) {
	timer := NewSuperTimer(5)
	for i:=0;i<5;i++{
		go push(timer, "worker" + strconv.Itoa(i))
	}
	time.Sleep(100)
}