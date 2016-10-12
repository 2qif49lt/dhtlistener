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
	len  int
}

func newSizeHashId(size int) *hashid {
	if size > hash_size*8 {
		panic("size is  bigger than 160")
	}
	ret := &hashid{}
	ret.len = size

	return ret
}

// newHashId takes string parameter and returns a hashid point
func newHashId(data string) *hashid {
	return newHashIdFromBytes([]byte(data))
}

// newHashId takes byte slice parameter and returns a hashid point
func newHashIdFromBytes(data []byte) *hashid {
	if len(data) > hash_size {
		panic("data length is bigger than 20")
	}

	id := &hashid{}
	copy(id.data[:], data)
	id.len = len(data) * 8

	return id
}

// Bit returns the idx-th position's bit, 0 or 1
func (h *hashid) Bit(idx int) int {
	if idx >= h.len {
		panic("index is out of range")
	}

	byteIdx, bitIdx := idx/8, idx%8
	return int((h.data[byteIdx] >> uint(8-bitIdx-1)) & 0x1)
}

// set sets the idx-th postions's bit to bit value
func (h *hashid) set(idx, bit int) {
	if idx >= h.len {
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
	if h.len != rhs.len {
		panic("size not the same")
	}

	ret := newSizeHashId(h.len)

	for k, _ := range h.data {
		ret.data[k] = h.data[k] ^ rhs.data[k]
	}

	return ret
}

// Compare compares the prefixLen-prefix of two bitmap.
//   - If bitmap.data[:prefixLen] < other.data[:prefixLen], return -1.
//   - If bitmap.data[:prefixLen] > other.data[:prefixLen], return 1.
//   - Otherwise return 0.
func (h *hashid) Compare(other *hashid, prefixLen int) int {
	if prefixLen > h.len || prefixLen > other.len {
		panic("index out of range")
	}
	div, mod := prefixLen/8, prefixLen%8

	for i := 0; i < div; i++ {
		if h.data[i] > other.data[i] {
			return 1
		} else if h.data[i] < other.data[i] {
			return -1
		}
	}

	for i := div * 8; i < div*8+mod; i++ {
		bit1, bit2 := h.Bit(i), other.Bit(i)
		if bit1 > bit2 {
			return 1
		} else if bit1 < bit2 {
			return -1
		}
	}

	return 0
}

func (h *hashid) String() string {
	div, mod := h.len/8, h.len%8
	var arr []string
	if mod > 0 {
		arr = make([]string, div+1)
	} else {
		arr = make([]string, div)
	}

	for k := 0; k != div; k++ {
		v := h.data[k]
		arr[k] = fmt.Sprintf("%08b", v)
	}
	if mod > 0 {
		modbin := ""
		for k := div * 8; k != div*8+mod; k++ {
			modbin += fmt.Sprintf("%d", h.Bit(k))
		}
		arr[div] = modbin
	}

	return strings.Join(arr[:], " ")
}

func (h *hashid) RawString() string {
	return string(h.data[:])
}

func (h *hashid) PrefixLen() int {
	for idx := 0; idx != h.len; idx++ {
		for idxbit := 0; idxbit != 8; idxbit++ {
			if h.data[idx]&(0x1<<(7-idxbit)) != 0 {
				return idx*8 + idxbit
			}
		}
	}

	return h.len*8 - 1
}
