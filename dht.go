package dhtlistener

import (
	"net"
	"strings"
)

type DHT struct {
	K    int
	me   *node
	addr string
	conn *net.UDPConn

	rt     *routetable
	peers  *peersManager
	tokens *tokenMgr
}

func NewDht(addr string) *DHT {
	var err error = nil

	var me *node = nil
	var udp_conn *net.UDPConn = nil

	if strings.Contains(addr, ":") {
		udp_addr, ok := net.ResolveUDPAddr("udp", addr)
		if ok == false {
			return nil
		}
		me = newRandomNodeFromUdpAddr(udp_addr)

		udp_conn, err = net.ListenUDP("udp", udp_addr)
		if err != nil {
			return nil
		}

	} else {
		addr += ":0"
		udp_addr, ok := net.ResolveUDPAddr("udp", addr)
		if ok == false {
			return nil
		}
		udp_conn, err = net.ListenUDP("udp", udp_addr)
		if err != nil {
			return nil
		}
		me = newRandomNodeFromUdpAddr(udp_conn.LocalAddr().(*net.UDPAddr))
	}

	ret := &DHT{
		k:      8,
		me:     me,
		addr:   addr,
		conn:   udp_conn,
		tokens: newTokenMgr(),
	}

	return ret
}

func (dht *DHT) init() {
	dht.rt = newRouteTable(dht)
	dht.peers = newPeersManager(dht)
}
