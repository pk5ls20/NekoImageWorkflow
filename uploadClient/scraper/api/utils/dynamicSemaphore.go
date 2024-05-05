package utils

import (
	"context"
	"sync"
)

type DynamicSemaphore struct {
	SetVal     int
	CurrentVal int
	Mutex      sync.Mutex
	Cond       *sync.Cond
}

func (ds *DynamicSemaphore) Acquire(ctx context.Context) error {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	for ds.CurrentVal >= ds.SetVal {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		ds.Cond.Wait()
	}
	ds.CurrentVal++
	return nil
}

func (ds *DynamicSemaphore) Release() {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	ds.CurrentVal--
	if ds.CurrentVal < ds.SetVal {
		ds.Cond.Signal()
	}
}

func (ds *DynamicSemaphore) AdjustSize(newSize int) {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	oldSize := ds.SetVal
	ds.SetVal = newSize
	if newSize > oldSize {
		for i := 0; i < (newSize - oldSize); i++ {
			ds.Cond.Signal()
		}
	}
}
