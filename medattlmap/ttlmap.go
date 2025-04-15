package medattlmap

// This is like the primitive replacement of REDIS. Basically we want to
// store key-value and automatically expires after X amount of time.
// It uses ticker that scan all for map for expiration and delete it == NOT EFFICIENT
// Also uses sync.Map

// Usage
// ttlMap := NewTTLMap(5 * time.Second) // Items expire after 5 seconds

// Store value
// ttlMap.Put("key1", "value1")

// Load value before expiration
// if val, ok := ttlMap.Get("key1"); ok {
//     println("Loaded:", val.(string)) // Output: Loaded: value1
// }

// Creating an instance of MyStruct
// myData := MyStruct{Name: "Example", Value: 42}

// Storing MyStruct in the TTL map
// ttlMap.Put("key1", myData)

// Load value before expiration
// if val, ok := ttlMap.Get("key1"); ok {
// Type casting to (MyStruct)
// 		fmt.Printf("Loaded: %+v\n", val.(MyStruct)) // Output: Loaded: {Name:Example Value:42}
// }

// Wait for 6 seconds to let the item expire
// time.Sleep(6 * time.Second)

// Attempt to load the expired key
// if _, ok := ttlMap.Get("key1"); !ok {
//     println("Key not found") // Output: Key not found
// }

// Close the TTLMap or leave it open, or in main put it in defer!
// ttlMap.Stop() // Stop the cleanup goroutine when done

import (
	"sync"
	"time"
)

const (
	DEFAULT_TICKER_TTL time.Duration = 5 * time.Second
	DEFAULT_TTL        time.Duration = 5 * time.Minute
)

// Value can be anything at this point, can be struct as well...
type item struct {
	value      interface{}
	expiration int64 // Unix timestamp for expiration
}

type TTLMap struct {
	m   sync.Map
	ttl time.Duration
	// tickerTTL time.Duration
	ticker *time.Ticker // global checker for all in this map
	stop   chan struct{}
}

// NewTTLMap creates a new TTLMap with the specified time-to-live
// if called with 0 and 0 then it's set to DEFAULT const above
// the tickttl is the value for "cron" checked, every tick we will check all
// the map for expiration
func NewTTLMap(ttl, tickttl time.Duration) *TTLMap {
	mandatory := ttl
	optional := tickttl
	if tickttl == 0 {
		optional = DEFAULT_TICKER_TTL
	}
	if ttl == 0 {
		mandatory = DEFAULT_TTL
	}
	t := &TTLMap{
		ttl: mandatory,
		// tickerTTL: optional,
		stop:   make(chan struct{}),
		ticker: time.NewTicker(optional), // Cleanup every second
	}

	go t.cleanup() // Start the cleanup goroutine
	return t
}

// Put adds or updates an item in the map with an expiration time.
func (t *TTLMap) Put(key string, ttl time.Duration, value interface{}) {
	optional := ttl
	if ttl == 0 {
		optional = t.ttl
	}
	expiration := time.Now().Add(optional).Unix()
	t.m.Store(key, &item{value: value, expiration: expiration})
}

// Get how many entries are in the ttlMap
func (t *TTLMap) Map() map[string]interface{} {
	vals := make(map[string]interface{})
	t.m.Range(func(k, v interface{}) bool {
		vals[k.(string)] = v
		return true
	})
	return vals
}

// Get how many entries are in the ttlMap
func (t *TTLMap) Len() int {
	var i int
	t.m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

// Get retrieves an item from the map if it exists and is not expired.
func (t *TTLMap) Get(key string) (interface{}, bool) {
	if v, ok := t.m.Load(key); ok {
		it := v.(*item)
		if time.Now().Unix() < it.expiration {
			return it.value, true
		}
		t.Delete(key) // Delete if expired
	}
	return nil, false
}

// Delete removes an item from the map.
func (t *TTLMap) Delete(key string) {
	t.m.Delete(key)
}

// Cleanup periodically removes expired items from the map.
// This is the ticker for checking expiration of the map
func (t *TTLMap) cleanup() {
	for {
		select {
		case <-t.ticker.C:
			now := time.Now().Unix()
			t.m.Range(func(key, value interface{}) bool {
				it := value.(*item)
				if now >= it.expiration {
					t.Delete(key.(string))
				}
				return true // continue iteration
			})
		case <-t.stop:
			t.ticker.Stop()
			return
		}
	}
}

// Stop stops the cleanup goroutine.
func (t *TTLMap) Stop() {
	close(t.stop)
}
