package utils

import (
	"errors"
	"sync"
)

type IDLock struct {
	lockMap sync.Map
}

type lockEntry struct {
	ch chan struct{}
}

// NewIDLock creates a new IDLock instance
// TODO: add TTL and optimise idLock impl
func NewIDLock() *IDLock {
	return &IDLock{}
}

func (l *IDLock) Lock(id string) {
	entry, loaded := l.lockMap.LoadOrStore(id, &lockEntry{
		ch: make(chan struct{}, 1),
	})
	e := entry.(*lockEntry)
	if !loaded {
		e.ch <- struct{}{}
	}
	<-e.ch
}

func (l *IDLock) Unlock(id string) error {
	if val, ok := l.lockMap.Load(id); ok {
		e := val.(*lockEntry)
		select {
		case e.ch <- struct{}{}:
		default:
			return errors.New("channel is full")
		}
	}
	return errors.New("no lock found for id")
}
