package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type IDLock struct {
	lockMap sync.Map
}

type lockEntry struct {
	ch chan struct{}
}

var (
	once     sync.Once
	instance *IDLock
)

// NewIDLock creates a new IDLock instance
// TODO: add TTL and optimise idLock impl
func NewIDLock() *IDLock {
	once.Do(func() {
		instance = &IDLock{}
	})
	return instance
}

func (l *IDLock) Lock(id string) {
	logrus.Debug("Start to Locking id: ", id)
	entry, loaded := l.lockMap.LoadOrStore(id, &lockEntry{
		ch: make(chan struct{}, 1),
	})
	e := entry.(*lockEntry)
	if !loaded {
		logrus.Debugf("Making new chan for id %s", id)
		e.ch <- struct{}{}
	}
	<-e.ch
	logrus.Debugf("Locked id: %s", id)
}

func (l *IDLock) Unlock(id string) error {
	if val, ok := l.lockMap.Load(id); ok {
		e := val.(*lockEntry)
		select {
		case e.ch <- struct{}{}:
			logrus.Debug("Unlocked id: ", id)
			return nil
		default:
			return fmt.Errorf("id %s channel is full", id)
		}
	}
	return fmt.Errorf("no lock found for id %s", id)
}
