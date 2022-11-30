package cache

import (
	"encoding/json"
	"io"
	"sync"
)

// key of the memory cache
type Key string

// value contained in memory
type Value interface{}

type Cache struct {
	values map[Key]Value
	lock   sync.RWMutex
}

type CacheJSON struct {
	Values map[Key]Value `json:"cache"`
}

const Threaded = true

func (c *Cache) Asleep(w io.Writer) error {
	payload := &CacheJSON{Values: c.values}
	return json.NewEncoder(w).Encode(payload)
}

func (c *Cache) Awake(r io.Reader) error {
	var payload CacheJSON
	err := json.NewDecoder(r).Decode(&payload)
	if err == nil {
		c.values = payload.Values
	}
	return err
}

func (c *Cache) Get(k Key) (Value, bool) {
	if Threaded {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}
	value, ok := c.values[k]
	if !ok {
		return nil, false
	}
	return value, true
}

func (c *Cache) Set(k Key, value Value) {
	if Threaded {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	c.values[k] = value
}

func (c *Cache) Remove(k Key) {
	if Threaded {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	delete(c.values, k)
}

func New() *Cache {
	cache := &Cache{
		values: make(map[Key]Value),
	}
	return cache
}
