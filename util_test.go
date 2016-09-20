package dhtlistener

import (
	"testing"
)

func TestAtoi64(t *testing.T) {
	cases := []struct {
		in  []byte
		out uint64
	}{
		{[]byte{0}, 0},
		{[]byte{1}, 1},
		{[]byte{1, 0}, 256},
		{[]byte{86, 113}, 22129},
		{[]byte{49, 50, 51}, 3224115},
	}

	for _, c := range cases {
		if Atoi64(c.in) != c.out {
			t.Fatal(c.in)
		}
	}
}

func TestI64toA(t *testing.T) {
	cases := []struct {
		in  uint64
		out []byte
	}{
		{0, []byte{0}},
		{1, []byte{1}},
		{256, []byte{1, 0}},
		{22129, []byte{86, 113}},
		{3224115, []byte{49, 50, 51}},
	}

	for _, c := range cases {
		bt := I64toA(c.in)
		for k, v := range bt {
			if v != c.out[k] {
				t.Fatal(c.in, k, v)
			}
		}
	}
}
