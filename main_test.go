package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestInMemoryCache_GetOrSet(t *testing.T) {
	c := NewInMemoryCache()

	inform := make(chan map[string]string, 1000)

	var ops1 uint64
	var ops2 uint64

	k := "key"
	v := "value"

	var mutex = &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < 1000; i++ {
		go func() {
			key := strconv.FormatUint(atomic.AddUint64(&ops1, 1), 10)
			value := strconv.FormatUint(atomic.AddUint64(&ops2, 1), 10)
			keys := fmt.Sprintf("%v%v", k, key)
			values := fmt.Sprintf("%v%v", v, value)

			z := map[string]string{
				keys: values,
			}
			inform <- z
		}()
	}
	wg.Done()


	wg.Add(1)

	for i := 0; i < 1000; i++  {
		go func() {
			mp := <- inform
			for k, v := range mp {
				c.GetOrSet(k,
					func() Value {
						mutex.Lock()
						mp[k] = v
						mutex.Unlock()
						return v
					},
				)
			}
		}()
	}

	wg.Done()

	wg.Wait()
	if len(c.data) == 0 {
		t.Errorf("there are %v items in the cache", len(c.data))
	}
}
