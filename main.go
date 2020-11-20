package main

import (
	"sync"
)

type (
	Key   = string
	Value = string
)

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

func main() {

}

func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
	cache.dataMutex.RLock()
	value, found := cache.Get(key)
		cache.dataMutex.RUnlock()
		cache.dataMutex.Lock()
	if found {
		cache.dataMutex.RUnlock()
		return value
	} else {
		value = valueFn()
		cache.data[key] = value
		cache.dataMutex.Unlock()
		return value
	}
}

