package dhtlistener

import (
	"errors"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// GetProcAbsDir returns the current process's folder path
func GetProcAbsDir() (string, error) {
	abs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", nil
	}
	return filepath.Dir(abs), nil
}

// ErrStr returns error's string or "nil"
func ErrStr(e error) string {
	if e == nil {
		return "nil"
	}
	return e.Error()
}

// GetRandString returns a size length random string
func GetRandString(size int) string {
	ret := make([]byte, size)
	rand.Read(ret)
	return string(ret)
}

// Atoi64 converts bytes to integer
func Atoi64(b []byte) uint64 {
	size, ret := len(b), uint64(0)
	if size > 8 {
		panic("input length great than 8")
	}

	for k, v := range b {
		ret |= uint64(v) << (uint32(size-k-1) * 8)
	}
	return ret
}

// I64toA converts integer to bytes
func I64toA(i uint64) []byte {
	ret := make([]byte, 8)
	count := 0

	for i != 0 {
		tmp := 0xff & i
		i = i >> 8
		ret[8-count-1] = byte(tmp)
		count++
	}

	return ret[8-count:]
}

// decodeCompactIPPortInfo decodes compactIP-address/port info in BitTorrent
// DHT Protocol. It returns the ip and port number.
func decodeCompactIPPortInfo(info string) (ip net.IP, port int, err error) {
	if len(info) != 6 {
		err = errors.New("compact info should be 6-length long")
		return
	}

	ip = net.IPv4(info[0], info[1], info[2], info[3])
	port = int((uint16(info[4]) << 8) | uint16(info[5]))
	return
}

// encodeCompactIPPortInfo encodes an ip and a port number to
// compactIP-address/port info.
func encodeCompactIPPortInfo(ip net.IP, port int) (info string, err error) {
	if port > 65535 || port < 0 {
		err = errors.New(
			"port should be no greater than 65535 and no less than 0")
		return
	}

	p := I64toA(uint64(port))
	if len(p) < 2 {
		p = append(p, p[0])
		p[0] = 0
	}

	info = string(append(ip, p...))
	return
}
