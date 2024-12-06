package faenet

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// FaerieGathering keeps track of magical connections
type FaerieGathering struct {
	sync.RWMutex
	activeFaeries int32
	totalFaeries  int32
}

func (fg *FaerieGathering) Add(n int) {
	atomic.AddInt32(&fg.totalFaeries, int32(n))
}

func (fg *FaerieGathering) Done() {
	atomic.AddInt32(&fg.activeFaeries, -1)
}

func (fg *FaerieGathering) DoneAll() {
	for atomic.LoadInt32(&fg.activeFaeries) > 0 {
		fg.Done()
	}
}

func (fg *FaerieGathering) Wait() {
	for atomic.LoadInt32(&fg.activeFaeries) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
}

func (fg *FaerieGathering) SummonNewFaerie() int32 {
	return atomic.AddInt32(&fg.totalFaeries, 1)
}

func (fg *FaerieGathering) WakeFaerie() {
	atomic.AddInt32(&fg.activeFaeries, 1)
}

func (fg *FaerieGathering) SlumberFaerie() {
	atomic.AddInt32(&fg.activeFaeries, -1)
}

func (fg *FaerieGathering) WhisperMagicalStats() string {
	return fmt.Sprintf("[%d/%d]", atomic.LoadInt32(&fg.activeFaeries), atomic.LoadInt32(&fg.totalFaeries))
}
