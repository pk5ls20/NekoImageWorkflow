package utils

import "sync"

type LockValue[T any] struct {
	Value T
	Lock  *sync.Mutex
}

func (lv *LockValue[T]) Get() T {
	lv.Lock.Lock()
	defer lv.Lock.Unlock()
	return lv.Value
}

func (lv *LockValue[T]) Set(value T) {
	lv.Lock.Lock()
	defer lv.Lock.Unlock()
	lv.Value = value
}

type RWLockValue[T any] struct {
	Value T
	Lock  *sync.RWMutex
}

func (lv *RWLockValue[T]) Get() T {
	lv.Lock.RLock()
	defer lv.Lock.RUnlock()
	return lv.Value
}

func (lv *RWLockValue[T]) Set(value T) {
	lv.Lock.Lock()
	defer lv.Lock.Unlock()
	lv.Value = value
}
