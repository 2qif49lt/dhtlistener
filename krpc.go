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

// makeQuery returns a query-formed data.
func makeQuery(t, q string, a map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"t": t,
		"y": "q",
		"q": q,
		"a": a,
	}
}

// makeResponse returns a response-formed data.
func makeResponse(t string, r map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"t": t,
		"y": "r",
		"r": r,
	}
}

// makeError returns a err-formed data.
func makeError(t string, errCode int, errMsg string) map[string]interface{} {
	return map[string]interface{}{
		"t": t,
		"y": "e",
		"e": []interface{}{errCode, errMsg},
	}
}

func send(dhr *DHT, addr *net.UDPAddr, data map[string]interface{}) error {

}
