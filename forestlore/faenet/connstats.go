package faenet

import (
	"fmt"
	"sync/atomic"
)

// FaerieGathering keeps track of magical connections
type FaerieGathering struct {
	totalFaeries  int32
	activeFaeries int32
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
