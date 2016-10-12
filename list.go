package dhtlistener

import (
	"container/list"
	"sync"
)

type syncList struct {
	*sync.RWMutex
	lst *list.List
}

func newSyncList() *syncList {
	return &syncList{&sync.RWMutex{}, list.New()}
}

func (sl *syncList) Back() interface{} {
	sl.RLock()
	defer sl.RUnlock()

	if sl.lst.Len() == 0 {
		return nil
	}

	return sl.lst.Back().Value
}

func (sl *syncList) Front() interface{} {
	sl.RLock()
	defer sl.RUnlock()

	if sl.lst.Len() == 0 {
		return nil
	}

	return sl.lst.Front().Value
}

func (sl *syncList) InsertAfter(v interface{}, mark *list.Element) *list.Element {
	sl.Lock()
	defer sl.Unlock()

	return sl.lst.InsertAfter(v, mark)
}

func (sl *syncList) InsertBefore(v interface{}, mark *list.Element) *list.Element {
	sl.Lock()
	defer sl.Unlock()

	return sl.lst.InsertBefore(v, mark)
}

func (sl *syncList) PushBack(v interface{}) *list.Element {
	sl.Lock()
	defer sl.Unlock()

	return sl.lst.PushBack(v)
}

func (sl *syncList) PushFront(v interface{}) *list.Element {
	sl.Lock()
	defer sl.Unlock()

	return sl.lst.PushFront(v)
}

func (sl *syncList) Remove(e *list.Element) interface{} {
	sl.Lock()
	defer sl.Unlock()

	return sl.lst.Remove(e)
}

func (sl *syncList) Clear() {
	sl.Lock()
	defer sl.Unlock()

	sl.lst.Init()
}

func (sl *syncList) Len() int {
	sl.RLock()
	defer sl.RUnlock()

	return sl.lst.Len()
}

func (sl *syncList) Foreach(f func(interface{}) bool) {
	sl.Lock()
	defer sl.Unlock()

	for e := sl.lst.Front(); e != nil; e = e.Next() {
		if f(e.Value) == false {
			break
		}
	}
}

func (sl *syncList) Iter() <-chan *list.Element {
	ch := make(chan *list.Element)

	go func() {
		sl.RLock()
		defer sl.RUnlock()

		for e := sl.lst.Front(); e != nil; e = e.Next() {
			ch <- e
		}
		close(ch)
	}()

	return ch
}

func (sl *syncList) Has(val interface{}) *list.Element {
	sl.RLock()
	defer sl.RUnlock()

	for e := sl.lst.Front(); e != nil; e = e.Next() {
		if e.Value == val {
			return e
		}
	}
	return nil
}
