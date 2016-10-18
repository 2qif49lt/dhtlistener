package dhtlistener

import (
	"net"
	"sync"
)

// Peer represents a peer contact.
type Peer struct {
	IP    net.IP
	Port  int
	token string
}

// newPeer returns a new peer pointer.
func newPeer(ip net.IP, port int, token string) *Peer {
	return &Peer{
		IP:    ip,
		Port:  port,
		token: token,
	}
}

// newPeerFromCompactIPPortInfo create a peer pointer by compact ip/port info.
func newPeerFromCompactIPPortInfo(compactInfo, token string) (*Peer, error) {
	ip, port, err := decodeCompactIPPortInfo(compactInfo)
	if err != nil {
		return nil, err
	}

	return newPeer(ip, port, token), nil
}

// CompactIPPortInfo returns "Compact node info".
// See http://www.bittorrent.org/beps/bep_0005.html.
func (p *Peer) CompactIPPortInfo() string {
	info, _ := encodeCompactIPPortInfo(p.IP, p.Port)
	return info
}

// peersManager represents a proxy that manipulates peers.
type peersManager struct {
	sync.RWMutex
	table *syncMap // hashinfo:peer
	dht   *DHT
}

// newPeersManager returns a new peersManager.
func newPeersManager(dht *DHT) *peersManager {
	return &peersManager{
		table: newsyncMap(),
		dht:   dht,
	}
}

// Insert adds a peer into peersManager.
func (pm *peersManager) Insert(infoHash string, peer *Peer) {
	pm.Lock()
	if _, ok := pm.table.Get(infoHash); !ok {
		pm.table.Set(infoHash, newSyncList())
	}
	pm.Unlock()

	v, _ := pm.table.Get(infoHash)
	queue := v.(*syncList)

	queue.PushBack(peer)
	if queue.Len() > pm.dht.K {
		queue.RemoveFront()
	}
}

// GetPeers returns size-length peers who announces having infoHash.
func (pm *peersManager) GetPeers(infoHash string, size int) []*Peer {
	peers := make([]*Peer, 0, size)

	v, ok := pm.table.Get(infoHash)
	if !ok {
		return peers
	}

	for e := range v.(*syncList).Iter() {
		peers = append(peers, e.Value.(*Peer))
	}

	if len(peers) > size {
		peers = peers[len(peers)-size:]
	}
	return peers
}
