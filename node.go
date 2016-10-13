package dhtlistener

import (
	"errors"
	"net"
	"strings"
	"time"
)

// node represents a DHT node.
type node struct {
	id             *hashid
	addr           *net.UDPAddr
	lastActiveTime time.Time
}

// newNode returns a node pointer.
func newNode(id, network, address string) (*node, error) {
	if len(id) != 20 {
		return nil, errors.New("node id should be a 20-length string")
	}

	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}

	return &node{newHashId(id), addr, time.Now()}, nil
}

func newRandomNodeFromUdpAddr(addr *net.UDPAddr) *node {
	return &node{newHashId(GetRandString(20)), addr, time.Now()}
}

// newNodeFromCompactInfo parses compactNodeInfo and returns a node pointer.
func newNodeFromCompactInfo(
	compactNodeInfo string, network string) (*node, error) {

	if len(compactNodeInfo) != 26 {
		return nil, errors.New("compactNodeInfo should be a 26-length string")
	}

	id := compactNodeInfo[:20]
	ip, port, _ := decodeCompactIPPortInfo(compactNodeInfo[20:])

	return newNode(id, network, genAddress(ip.String(), port))
}

// CompactIPPortInfo returns "Compact IP-address/port info".
// See http://www.bittorrent.org/beps/bep_0005.html.
func (node *node) CompactIPPortInfo() string {
	info, _ := encodeCompactIPPortInfo(node.addr.IP, node.addr.Port)
	return info
}

// CompactNodeInfo returns "Compact node info".
// See http://www.bittorrent.org/beps/bep_0005.html.
func (node *node) CompactNodeInfo() string {
	return strings.Join([]string{
		node.id.RawString(), node.CompactIPPortInfo(),
	}, "")
}

type sortNodeByTime []*node

func (st sortNodeByTime) Len() int {
	return len(st)
}

func (st sortNodeByTime) Swap(i, j int) {
	st[i], st[j] = st[j], st[i]
}
func (st sortNodeByTime) Less(i, j int) bool {
	return st[i].lastActiveTime.Before(st[j].lastActiveTime)
}

type sortNodeById []*node

func (st sortNodeById) Len() int {
	return len(st)
}

func (st sortNodeById) Swap(i, j int) {
	st[i], st[j] = st[j], st[i]
}

func (st sortNodeById) Less(i, j int) bool {
	return st[i].id.Compare(st[j].id, st[i].id.len) == -1
}
