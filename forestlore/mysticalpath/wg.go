package mysticalpath

import (
	"sync"
	"sync/atomic"
)

type faerieGathering struct {
	innerCircle sync.WaitGroup
	faerieCount int32
}

func (fg *faerieGathering) SummonFaeries(n int) {
	atomic.AddInt32(&fg.faerieCount, int32(n))
	fg.innerCircle.Add(n)
}

func (fg *faerieGathering) FaerieDeparted() {
	if count := atomic.LoadInt32(&fg.faerieCount); count > 0 && atomic.CompareAndSwapInt32(&fg.faerieCount, count, count-1) {
		fg.innerCircle.Done()
	}
}

func (fg *faerieGathering) AllFaeriesDeparted() {
	for atomic.LoadInt32(&fg.faerieCount) > 0 {
		fg.FaerieDeparted()
	}
}

func (fg *faerieGathering) AwaitFaerieGathering() {
	fg.innerCircle.Wait()
}
