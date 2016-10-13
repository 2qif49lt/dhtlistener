package dhtlistener

import (
	"container/list"
	"sync"
	"time"
)

type keylist struct {
	*sync.RWMutex
	*syncList
	keyMap map[interface{}]*list.Element
}

func newKeyList() *keylist {
	return &keylist{
		&sync.RWMutex{},
		newSyncList(),
		make(map[interface{}]*list.Element),
	}
}

func (kl *keylist) Get(key interface{}) (interface{}, bool) {
	kl.RLock()
	defer kl.RUnlock()

	e, ok := kl.keyMap[key]
	if ok {
		return e.Value, ok
	}
	return nil, false
}

func (kl *keylist) Has(key interface{}) bool {
	_, exist := kl.Get(key)

	return exist
}

func (kl *keylist) Push(key, val interface{}) {
	kl.Lock()
	defer kl.Unlock()

	e, exist := kl.keyMap[key]
	if exist {
		kl.Remove(e)
	}

	e = kl.PushBack(val)
	kl.keyMap[key] = e

}

func (kl *keylist) Remove(key interface{}) interface{} {
	kl.Lock()
	defer kl.Unlock()

	e, exist := kl.keyMap[key]
	if exist {
		kl.Remove(e)
		delete(kl.keyMap, key)
		return e.Value
	}
	return nil
}

func (kl *keylist) Clear() {
	kl.Lock()
	defer kl.Unlock()

	kl.syncList.Clear()
	kl.keyMap = make(map[interface{}]*list.Element)
}
