package dhtlistener

import (
	"sort"
)

type routetable struct {
	dht     *DHT
	buckets [hash_size * 8]*keylist // rawstring:*node
}

func newRouteTable(dht *DHT) *routetable {
	ret := &routetable{
		dht: dht,
	}

	for idx := 0; idx != len(ret.buckets); idx++ {
		ret.buckets[idx] = newKeyList()
	}

	return ret
}

func (rt *routetable) Insert(n *node) bool {
	prefix_len := n.id.Xor(rt.dht.me.id).PrefixLen()
	bucket := rt.buckets[prefix_len]

	if bucket.Has(n.id.RawString()) {
		bucket.Remove(n.id.RawString())
		bucket.Push(n.id.RawString(), n)
		return false
	}
	if bucket.Len() < rt.dht.K {
		bucket.Push(n.id.RawString(), n)
		return true
	} else {
		go func(l *keylist) {
			// ping bucket
		}(bucket)
	}
	return false
}

func (rt *routetable) getNode(h string) *node {
	for _, v := range rt.buckets {
		if n, exist := v.Get(h); exist {
			if n.(*node).id.RawString() == h {
				return n.(*node)
			} else {
				return nil
			}
		}
	}
	return nil
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
	if len(ret) > size {
		ret = ret[:size]
	}
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

func (rt *routetable) Fresh() {
	for idx, bucket := range rt.buckets {

		bucket.Foreach(func(it interface{}) bool {
			n := it.(*node)
			newid := newSubRandHashIdFromIdSize(rt.dht.me.id, idx)

			_, _ = n, newid
			// go find node
			return true
		})

	}
}

func (rt *routetable) Len() int {
	ret := 0
	for _, bucket := range rt.buckets {
		ret += bucket.Len()
	}
	return ret
}
