package queue

import (
	"testing"
	"time"
)

func TestBinaryHead(t *testing.T) {

	// new
	b := NewBinaryHead()
	// test data
	data := map[int]int{
		1:  1,
		2:  2,
		3:  3,
		4:  4,
		5:  5,
		6:  6,
		7:  7,
		8:  8,
		9:  9,
		10: 10,
	}
	// push
	log.Info("push")
	for _, val := range data {
		log.Warning("push val: %d", val)
		b.Push(val, val)
	}
	// pop
	log.Info("pop")

	for range data {
		val, _ := b.Pop()
		log.Error("pop val: %d", val.Val)
	}

	v, ok := b.Pop()
	log.Error("latest %v:%v", ok, v)

	log.Info("done")
	time.Sleep(time.Second)

}
