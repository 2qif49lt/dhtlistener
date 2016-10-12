package dhtlistener

import (
	"sync"
	"time"
)

type routetable struct {
	*sync.RWMutex
	dht     *DHT
	me      *node
	buckets [hash_size * 8]*newKeyList // rawstring:*node
}

func newRouteTable(me *node) *routetable {
	ret := &routetable{
		&sync.RWMutex,
		me,
	}

	for idx := 0; idx != len(ret.buckets); idx++ {
		ret.buckets[idx] = newKeyList()
	}

	return ret
}

func (rt *routetable) update(n *node) {
	rt.Lock()
	defer rt.Unlock()

	prefix_len := n.id.Xor(rt.me.id).PrefixLen()
	lst := rt.buckets[prefix_len]

	if lst.Has(n.id.RawString()) {
		lst.Remove(n.id.RawString())
		lst.Push(n.id.RawString(), n)
	} else {
		if lst.Len() < rt.dht.K {
			lst.Push(n.id.RawString(), n)
		} else {
			go func(l *newKeyList) {
				// ping bucket
			}(lst)
		}
	}
}

func (rt *routetable) FindClosest(tar *hashid, size int) []*node {
	rt.RLock()
	defer rt.RUnlock()

	ret := make([]*node, 0)

	bucket_num := tar.Xor(rt.me.id).PrefixLen()
	bucket := rt.buckets[bucket_num]

	bucket.Foreach(func(it interface{}) bool {
		n := it.(*node)
		ret = append(ret, n)
		return true
	})

	for i := 1; (bucket_num-1 >= 0 || bucket_num+1 < len(rt.buckets)) && len(ret) < size; i++ {

	}
}
