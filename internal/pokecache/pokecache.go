package pokecache

import (
	"sync"
	"time"
)

//timer := time.NewTicker(time.Second) //will need this for reaping interface
//struct within Cache struct
type cacheEntry struct {
	createdAt    time.Time
	Val        []byte
}
//main Cache struct
type Cache struct {
	mu	    sync.Mutex 
	Data    map[string]cacheEntry
}

//interface used in NewCache() call to remove old entries
func (c *Cache) reapLoop(interval time.Duration) {
	//delete(map, key)
	var start time.Time
	var current time.Time
	var duration time.Duration
	
	reapTicker := time.NewTicker(interval)
	for range reapTicker.C {
		c.mu.Lock()
		for key, val := range c.Data {
		start = val.createdAt
		current = time.Now()
		duration = current.Sub(start)
		if duration > interval {
			//c.mu.Lock()
			delete(c.Data, key)
			//reapticker does not end, using defer for the unlock causes a freeze, MOVED to outer loop before reading occurs
			//c.mu.Unlock()
			}
		} // end Cache data range loop
		c.mu.Unlock()
	} //end reapTicker loop
	
	return
}
//create a new referenced cache
func NewCache(interval time.Duration) *Cache {
	
	cacheOrigin := Cache{
		Data:    map[string]cacheEntry{
			"default": {
			createdAt:	time.Now(),
			Val:		make([]byte, 1),
			},
		},					
	}
	cacheOut := &cacheOrigin
	_, ok := cacheOut.Data["default"]
	if !ok {
		return cacheOut
	}
	go cacheOut.reapLoop(interval)

	return cacheOut
} //end func

//Cache interface, Add: adds new entries to the cache
func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	c.Data[key] = cacheEntry {
		createdAt:   time.Now(),
		Val:		 value,
	}
	defer c.mu.Unlock()

	return
}//end func

//Cache interface, Get: gets entry data from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.Data[key]
	//if ok, data acquisition successful
	if ok {
		return entry.Val, true
	}
	//data acquisition has failed
	fail := make([]byte, 1)
	return fail, false
}//end func

