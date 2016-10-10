package dhtlistener

import (
	"fmt"
	"testing"
)

func TestHashIdSimple(t *testing.T) {
	data := "0123456789abcdefghij"

	id := newHashId(data)
	id.len = 51
	fmt.Println(id)
	fmt.Println(id.RawString())

	data = "0123456789"
	id = newHashId(data)
	id.len += 1
	fmt.Println(id)
	fmt.Println(id.RawString())
}

func TestHashIdBit(t *testing.T) {
	data := "0123456789abcdefghij"
	/*
			   00110000 00110001 00110010 00110011 00110100
		       00110101 00110110 00110111 00111000 00111001
			   01100001 01100010 01100011 01100100 01100101
		       01100110 01100111 01101000 01101001 01101010
	*/

	cases := []struct {
		in  int
		out int
	}{
		{0, 0},
		{7, 0},
		{24, 0},
		{81, 1},
		{158, 1},
	}

	id := newHashId(data)
	for k, v := range cases {
		if id.Bit(v.in) != v.out {
			t.Fatal(k, id.Bit(v.in), v.out)
		}
	}
}

func TestHashIdSet(t *testing.T) {
	data := "0123456789abcdefghij"

	cases := []struct {
		in  int
		out string
	}{
		{7, "1123456789abcdefghij"},
		{23, "0133456789abcdefghij"},
		{87, "0123456789abcdefghij"},
		{159, "0123456789abcdefghik"},
	}

	for k, v := range cases {
		id := newHashId(data)

		id.Set(v.in)
		if id.RawString() != v.out {
			t.Fatal(k, id.RawString(), v.out)
		}

	}
}

func TestHashIdXor(t *testing.T) {
	data := "0123456789abcdefghij"
	inverseHashId := newHashId(data)

	for i := 80; i != 160; i++ {
		inverseHashId.Inverse(i)
	}

	/*
			   00110000 00110001 00110010 00110011 00110100
		       00110101 00110110 00110111 00111000 00111001
			   10011110 10011101 10011100 10011011 10011010
		       10011001 10011000 10010111 10010110 10010101
	*/

	halfZerohalfOne := newSizeHashId(160)
	for i := 80; i != 160; i++ {
		halfZerohalfOne.Set(i)
	}
	/*
			   00000000 00000000 00000000 00000000 00000000
		       00000000 00000000 00000000 00000000 00000000
			   11111111 11111111 11111111 11111111 11111111
		       11111111 11111111 11111111 11111111 11111111
	*/

	cases := []struct {
		in  string
		out string
	}{
		{"0123456789abcdefghij", string(make([]byte, 20))},
		{inverseHashId.RawString(), halfZerohalfOne.RawString()},
	}

	for k, v := range cases {
		id := newHashId(data)
		rhs := newHashId(v.in)
		if id.Xor(rhs).RawString() != v.out {
			t.Fatal(k, v)
		}
	}
}

func TestHashIdCompare(t *testing.T) {
	halfZerohalfOne := newSizeHashId(160)
	for i := 80; i != 160; i++ {
		halfZerohalfOne.Set(i)
	}

	tmp := *halfZerohalfOne
	other := &tmp
	other.UnSet(100)

	rst := halfZerohalfOne.Compare(other, 120)
	if rst != 1 {
		t.Fatal(rst, halfZerohalfOne, other)
	}
}
