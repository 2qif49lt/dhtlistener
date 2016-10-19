package dhtlistener

import (
	"encoding/hex"
	"net"
	"strings"
	"time"
)

type DHT struct {
	K              int
	me             *node
	addr           string
	conn           *net.UDPConn
	Try            int
	EntranceAddrs  []string
	packets        chan packet
	works          chan struct{}
	rt             *routetable
	peers          *peersManager
	transacts      *transactionManager
	tokens         *tokenMgr
	OnGetPeers     func(string, string, int)
	OnAnnouncePeer func(string, string, int)
}

func NewDht(addr string) *DHT {
	var err error = nil

	var me *node = nil
	var udp_conn *net.UDPConn = nil

	if strings.Contains(addr, ":") {
		udp_addr, ok := net.ResolveUDPAddr("udp", addr)
		if ok != nil {
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
		if ok != nil {
			return nil
		}
		udp_conn, err = net.ListenUDP("udp", udp_addr)
		if err != nil {
			return nil
		}
		me = newRandomNodeFromUdpAddr(udp_conn.LocalAddr().(*net.UDPAddr))
	}

	ret := &DHT{
		K:    8,
		me:   me,
		addr: addr,
		conn: udp_conn,
		Try:  2,
		EntranceAddrs: []string{
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
		},
		packets:        make(chan packet, 1024),
		works:          make(chan struct{}, 100),
		tokens:         newTokenMgr(),
		OnGetPeers:     nil,
		OnAnnouncePeer: nil,
	}

	return ret
}

func (dht *DHT) init() {
	dht.rt = newRouteTable(dht)
	dht.peers = newPeersManager(dht)
	dht.transacts = newTransactionManager(dht)
}

func (dht *DHT) srv() {
	go func() {
		buff := make([]byte, 8192)
		for {
			n, raddr, err := dht.conn.ReadFromUDP(buff)
			if err != nil {
				continue
			}

			dht.packets <- packet{buff[:n], raddr, time.Now()}
		}
	}()
}

func (dht *DHT) join() {
	for _, addr := range dht.EntranceAddrs {
		raddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			continue
		}
		dht.transacts.findNode(
			&node{addr: raddr},
			dht.me.id.RawString(),
		)
	}
}

func (dht *DHT) GetPeers(infoHash string) ([]*Peer, error) {
	if len(infoHash) == 40 {
		data, err := hex.DecodeString(infoHash)
		if err != nil {
			return nil, err
		}
		infoHash = string(data)
	}

	peers := dht.peers.GetPeers(infoHash, dht.K)
	if len(peers) != 0 {
		return peers, nil
	}

	ch := make(chan struct{})

	go func() {
		neighbors := dht.rt.FindClosestNode(newHashId(infoHash), dht.K)

		for _, no := range neighbors {
			dht.transacts.getPeers(no, infoHash)
		}

		i := 0
		for _ = range time.Tick(time.Second * 1) {
			i++
			peers = dht.peers.GetPeers(infoHash, dht.K)
			if len(peers) != 0 || i == 30 {
				break
			}
		}

		ch <- struct{}{}
	}()

	<-ch
	return peers, nil
}

func (dht *DHT) Run() {
	dht.init()
	dht.srv()
	dht.join()

	var pkt packet

	for {
		select {
		case pkt = <-dht.packets:
			handle(dht, pkt)
		case <-time.After(time.Second * 5):
			if dht.rt.Len() == 0 {
				dht.join()
			} else {
				go dht.rt.Fresh()
			}
		}
	}
}
