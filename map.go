package dhtlistener

import (
	"sync"
)

type mapItem struct {
	key interface{}
	val interface{}
}

// syncMap represents a goroutine-safe map.
type syncMap struct {
	*sync.RWMutex
	data map[interface{}]interface{}
}

// newsyncMap returns a syncMap pointer.
func newsyncMap() *syncMap {
	return &syncMap{
		RWMutex: &sync.RWMutex{},
		data:    make(map[interface{}]interface{}),
	}
}

// Get returns the value mapped to key.
func (smap *syncMap) Get(key interface{}) (val interface{}, ok bool) {
	smap.RLock()
	defer smap.RUnlock()

	val, ok = smap.data[key]
	return
}

// Has returns whether the syncMap contains the key.
func (smap *syncMap) Has(key interface{}) bool {
	_, ok := smap.Get(key)
	return ok
}

// Set sets pair {key: val}.
func (smap *syncMap) Set(key interface{}, val interface{}) {
	smap.Lock()
	defer smap.Unlock()

	smap.data[key] = val
}

// Delete deletes the key in the map.
func (smap *syncMap) Delete(key interface{}) {
	smap.Lock()
	defer smap.Unlock()

	delete(smap.data, key)
}

// DeleteMulti deletes keys in batch.
func (smap *syncMap) DeleteMulti(keys []interface{}) {
	smap.Lock()
	defer smap.Unlock()

	for _, key := range keys {
		delete(smap.data, key)
	}
}

// Clear resets the data.
func (smap *syncMap) Clear() {
	smap.Lock()
	defer smap.Unlock()

	smap.data = make(map[interface{}]interface{})
}

// Iter returns a chan which output all items.
func (smap *syncMap) Iter() <-chan mapItem {
	ch := make(chan mapItem)
	go func() {
		smap.RLock()
		for key, val := range smap.data {
			ch <- mapItem{
				key: key,
				val: val,
			}
		}
		smap.RUnlock()
		close(ch)
	}()
	return ch
}

// Len returns the length of syncMap.
func (smap *syncMap) Len() int {
	smap.RLock()
	defer smap.RUnlock()

	return len(smap.data)
}
