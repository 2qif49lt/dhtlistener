package dhtlistener

import (
	"net"
	"time"
)

const (
	// token‘s validity, unit sencond
	token_active_time = 60 * 10
)

// token used in getpeer and announce_peer
type token struct {
	data       string
	createTime int64
}

type tokenMgr struct {
	*syncMap
}

// newTokenMgr returns a new tokenManager.
func newTokenMgr() *tokenMgr {
	return &tokenMgr{
		syncMap: newsyncMap(),
	}
}

// getToken returns a token.
func (tm *tokenMgr) getToken(addr *net.UDPAddr) string {
	v, ok := tm.Get(addr.IP.String())
	tk, _ := v.(token)

	if !ok || time.Now().Unix()-tk.createTime > token_active_time {
		tk = token{
			data:       GetRandString(4),
			createTime: time.Now().Unix(),
		}

		tm.Set(addr.IP.String(), tk)
	}

	return tk.data
}

// clear removes expired tokens.
func (tm *tokenMgr) clearExpired() {
	for _ = range time.Tick(time.Minute * 3) {
		keys := make([]interface{}, 0, 100)

		for item := range tm.Iter() {
			if time.Now().Unix()-item.val.(token).createTime > token_active_time {
				keys = append(keys, item.key)
			}
		}

		tm.DeleteMulti(keys)
	}
}

// check returns whether the token is valid.
func (tm *tokenMgr) check(addr *net.UDPAddr, tokenString string) bool {
	key := addr.IP.String()
	v, ok := tm.Get(key)
	tk, _ := v.(token)

	if ok {
		tm.Delete(key)
	}

	return ok && tokenString == tk.data
}
