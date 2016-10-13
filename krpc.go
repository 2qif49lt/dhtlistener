package dhtlistener

import (
	"math"
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
	_, err := dht.conn.WriteToUDP([]byte(Encode(data)), addr)
	return err
}

// query represents the query data included queried node and query-formed data.
type query struct {
	tar  *node
	data map[string]interface{}
}

// transaction implements transaction.
type transaction struct {
	*query
	id       string
	response chan struct{}
}

type transactionManager struct {
	*sync.RWMutex
	transactions *syncMap
	index        *syncMap
	curTransId   uint64 // MaxInt32
	queryChan    chan *query
	dht          *DHT
}

func newTransactionManager(dht *DHT) *transactionManager {
	return &transactionManager{
		RWMutex:      &sync.RWMutex{},
		transactions: newsyncMap(),
		index:        newsyncMap(),
		queryChan:    make(chan *query, 1024),
		dht:          dht,
	}
}

// genTransID generates a transaction id and returns it.
func (tm *transactionManager) genTransID() string {
	tm.Lock()
	defer tm.Unlock()

	tm.curTransId = (tm.curTransId + 1) % math.MaxUInt32
	return string(I64toA(tm.curTransId))
}

func (tm *transactionManager) newTransaction(id string, q *query) *transaction {
	return &transaction{
		id:       id,
		query:    q,
		response: make(chan struct{}, tm.dht.Try+1),
	}
}
