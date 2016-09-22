package dhtlistener

import (
	"fmt"
	"strings"
)

const (
	hash_size = 20
)

// hashid holds information about node id
type hashid struct {
	data [hash_size]byte
}

// newHashId takes string parameter and returns a hashid point
func newHashId(data string) *hashid {
	return newHashIdFromBytes([]byte(data))
}

// newHashId takes byte slice parameter and returns a hashid point
func newHashIdFromBytes(data []byte) *hashid {
	if len(data) != hash_size {
		panic("data length is not 20")
	}

	id := &hashid{}
	copy(id.data[:], data)

	return id
}

// Bit returns the idx-th position's bit, 0 or 1
func (h *hashid) Bit(idx int) int {
	if idx >= 8*hash_size {
		panic("index is out of range")
	}

	byteIdx, bitIdx := idx/8, idx%8
	return int((h.data[byteIdx] >> uint(8-bitIdx-1)) & 0x1)
}

// set sets the idx-th postions's bit to bit value
func (h *hashid) set(idx, bit int) {
	if idx >= 8*hash_size {
		panic("index is out of range")
	}

	byteIdx, bitIdx := idx/8, idx%8
	opt := byte(1 << uint(8-bitIdx-1))
	if bit != 0 {
		h.data[byteIdx] |= opt
	} else {
		h.data[byteIdx] &= ^opt
	}
}

// Inverse inverse the idx-th position bit,returns the old value
func (h *hashid) Inverse(idx int) int {
	old := h.Bit(idx)
	if old == 0 {
		h.Set(idx)
	} else {
		h.UnSet(idx)
	}
	return old
}

// Set sets the idx-th postions's bit to 1
func (h *hashid) Set(idx int) {
	h.set(idx, 1)
}

// UnSet sets the idx-th postions's bit to 0
func (h *hashid) UnSet(idx int) {
	h.set(idx, 0)
}

// Xor returns xor value of two hashid
func (h *hashid) Xor(rhs *hashid) *hashid {
	ret := &hashid{}

	for k, _ := range h.data {
		ret.data[k] = h.data[k] ^ rhs.data[k]
	}

	return ret
}

func (h *hashid) String() string {
	arr := [20]string{}

	for k, v := range h.data {
		arr[k] = fmt.Sprintf("%08b", v)
	}

	return strings.Join(arr[:], " ")
}

func (h *hashid) RawString() string {
	return string(h.data[:])
}
