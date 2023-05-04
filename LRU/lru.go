package LRU

import (
	"container/list"
	"unsafe"
)

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	// max bytes that can be used.
	maxBytes int64
	// Bytes that have used
	nBytes int64
	// double linked list
	ll    *list.List
	cache map[string]*list.Element
	// optional and executed when an entry in purged.
	OnEvicted func(key string, value Value)
}

// Value uses Size() to count how many bytes it takes
type Value interface {
	Size() int
}

type entry struct {
	key   string
	value Value
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Size()) - int64(kv.value.Size())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(value.Size()) + int64(unsafe.Sizeof(key))
	}
	for c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(unsafe.Sizeof(kv.key)) + int64(kv.value.Size())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}
