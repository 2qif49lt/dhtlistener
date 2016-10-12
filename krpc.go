package dhtlistener

import (
	"net"
	"time"
)

const (
	pingNode     = "ping"
	findNode     = "find_node"
	getPeers     = "get_peers"
	announcePeer = "announce_peer"
)

const (
	genericError  = 201
	serverError   = 202
	protocolError = 203
	unknownError  = 204
)

// packet represents the information receive from udp.
type packet struct {
	data     []byte
	raddr    *net.UDPAddr
	recvTime time.Time
}
