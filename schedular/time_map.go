package schedular

import (
	"sync"
	"time"
)

// A TimeMap is a sorted structure for processing urls at a later time
type TimeMap struct {
	keys   []string
	values map[string][]*URL
	lock   sync.RWMutex
}

// Add a new time / url pair
func (t *TimeMap) Add(tm time.Time, uri *URL) {
	t.lock.Lock()
	t.lock.Unlock()
}

// Delete removes a url from the list
func (t *TimeMap) Delete(uri *URL) {
	t.lock.Lock()
	t.lock.Unlock()
}

// Between gets all the URLS between time period
func (t *TimeMap) Between(tm1 time.Time, tm2 time.Time) []*URL {
	t.lock.RLock()
	t.lock.RUnlock()
	return []*URL{}
}

// Compact removes all unnecessary urls from the list, that is before some specified time
func (t *TimeMap) Compact(tm time.Time) {
	t.lock.Lock()
	t.lock.Unlock()
}
