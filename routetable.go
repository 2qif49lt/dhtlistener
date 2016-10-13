package dhtlistener

import (
	"sort"
	"sync"
	"time"
)

type routetable struct {
	dht     *DHT
	buckets [hash_size * 8]*newKeyList // rawstring:*node
}

func newRouteTable(dht *DHT) *routetable {
	ret := &routetable{
		dht,
	}

	for idx := 0; idx != len(ret.buckets); idx++ {
		ret.buckets[idx] = newKeyList()
	}

	return ret
}

func (rt *routetable) update(n *node) {
	prefix_len := n.id.Xor(rt.dht.me.id).PrefixLen()
	bucket := rt.buckets[prefix_len]

	if bucket.Has(n.id.RawString()) {
		bucket.Remove(n.id.RawString())
		bucket.Push(n.id.RawString(), n)
	} else {
		if bucket.Len() < rt.dht.K {
			bucket.Push(n.id.RawString(), n)
		} else {
			go func(l *newKeyList) {
				// ping bucket
			}(bucket)
		}
	}
}

func (rt *routetable) FindClosestNode(tar *hashid, size int) []*node {
	ret := make([]*node, 0)

	bucket_num := tar.Xor(rt.dht.me.id).PrefixLen()
	bucket := rt.buckets[bucket_num]

	bucket.Foreach(func(it interface{}) bool {
		n := it.(*node)
		ret = append(ret, n)
		return true
	})

	for i := 1; (bucket_num-i >= 0 || bucket_num+i < len(rt.buckets)) && len(ret) < size; i++ {
		if bucket_num-i >= 0 {
			bucket = rt.buckets[bucket_num-i]
			bucket.Foreach(func(it interface{}) bool {
				n := it.(*node)
				ret = append(ret, n)
				return true
			})
		}
		if bucket_num+i < len(rt.buckets) {
			bucket = rt.buckets[bucket_num+i]
			bucket.Foreach(func(it interface{}) bool {
				n := it.(*node)
				ret = append(ret, n)
				return true
			})
		}
	}
	sort.Sort(sortNodeById(ret))

	return ret
}

func (rt *routetable) GetClosestNodeCompactInfo(tar *hashid, size int) []string {
	nodes := rt.FindClosestNode(tar, size)
	infos := make([]string, len(nodes))

	for k, v := range nodes {
		infos[k] = v.CompactNodeInfo()
	}
	return infos
}

func (rt *routetable) Remove(tar *hashid) {
	bucket_num := tar.Xor(rt.dht.me.id).PrefixLen()
	bucket := rt.buckets[bucket_num]

	bucket.Remove(tar.RawString())
}

func (rt *routetable) Flush() {
	for idx, bucket := range rt.buckets {

		bucket.Foreach(func(it interface{}) bool {
			n := it.(*node)
			newid := newSubRandHashIdFromIdSize(rt.dht.me.id, idx)

			// go find node
		})

	}
}
