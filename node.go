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
