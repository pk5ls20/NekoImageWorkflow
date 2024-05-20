package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

func TestBasicIDLock(t *testing.T) {
	lock := NewIDLock()
	id := "test-id"
	done := make(chan bool)
	// A. Test basic lock and unlock
	// 1. lock id , chan input = 1, after lock should not block
	go func() {
		lock.Lock(id)
		done <- true
	}()
	select {
	// done chan received, which means not block, expected
	case <-done:
	// done chan not received, which means block, not expected
	case <-time.After(1 * time.Second):
		t.Errorf("First lock timed out, it should not block")
	}
	// 2. lock id again, chan input = 2, after lock should block
	go func() {
		lock.Lock(id)
		done <- true
	}()
	select {
	// done chan received, which means not block, not expected
	case <-done:
		t.Errorf("Second lock did not block, it should block")
	// done chan not received, which means block, expected
	case <-time.After(1 * time.Second):
	}
	// 3. unlock id, chan input = 1, after unlock should not block
	if err := lock.Unlock(id); err != nil {
		t.Errorf("Unlock failed: %s", err)
	}
	// manually release done chan
	<-done
	// 4. lock id again, chan input = 2, after lock should block
	go func() {
		lock.Lock(id)
		logrus.Debug(11111)
		done <- true
	}()
	select {
	// done chan received, which means not block, not expected
	case <-done:
		t.Errorf("Second lock did not block, it should block")
	// done chan not received, which means block, expected
	case <-time.After(1 * time.Second):
	}
	// 5. unlock id x2
	if err := lock.Unlock(id); err != nil {
		t.Errorf("Unlock failed: %s", err)
	}
	if err := lock.Unlock(id); err != nil {
		t.Errorf("Unlock failed: %s", err)
	}
	// B. Test unlock non-existent id
	if err := lock.Unlock("non-existent-id"); err == nil {
		t.Errorf("Expected error when unlocking non-existent id, got none")
	}
}

func TestIDLockConcurrent(t *testing.T) {
	lock := NewIDLock()
	numGoroutines := 100
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			id := fmt.Sprintf("test-id-%d", i)
			lock.Lock(id)
			go lock.Lock(id)
			go lock.Lock(id)
			go lock.Lock(id)
			go func() {
				if err := lock.Unlock(id); err != nil {
					t.Errorf("Unlock failed for id %s: %s", id, err)
				}
				wg.Done()
			}()
		}(i)
	}
	wg.Wait()
}

func TestIDLockSameIDConcurrent(t *testing.T) {
	lock := NewIDLock()
	id := "test-id"
	numGoroutines := 100
	var wg sync.WaitGroup
	done := make(chan struct{})
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lock.Lock(id)
			err := lock.Unlock(id)
			if err != nil {
				t.Errorf("Unlock failed for id %s: %s", id, err)
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Errorf("Test timed out")
	}
}
