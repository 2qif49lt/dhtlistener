package dhtlistener

import (
	"math/rand"
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
