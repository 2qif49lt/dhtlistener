package dhtlistener

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	cases := []struct {
		in  int
		out string
	}{
		{1, "i1e"},
		{0, "i0e"},
		{42, "i42e"},
		{-42, "i-42e"},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := EncodeInt(cases[idx].in); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}
func TestEncodeString(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"spam", "4:spam"},
		{"hello,中国", "12:hello,中国"},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := EncodeString(cases[idx].in); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}

func TestEncodeSlice(t *testing.T) {
	cases := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{1, "spam"}, "li1e4:spame"},
		{[]interface{}{1, "spam", -1}, "li1e4:spami-1ee"},
		{[]interface{}{1, "spam", []int{1, 2}}, "li1e4:spamli1ei2eee"},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := EncodeSlice(cases[idx].in); err == nil {
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}

func TestEncodeMap(t *testing.T) {
	in := make(map[string]interface{})
	in["q"] = "ping"
	in["id"] = "identify"
	in["t"] = 123
	in["list"] = []string{"abc", "def"}

	out := "d2:id8:identify4:listl3:abc3:defe1:q4:ping1:ti123ee"
	if str, err := EncodeMap(in); err == nil {
		fmt.Println(str, out)
		if str != out {
			t.Fatal(str)
		}
	} else {
		t.Fatal(err)
	}
}

func TestEncodeStruct(t *testing.T) {
	in := struct {
		Q    string   `json:"q"`
		Id   string   `json:"id"`
		T    int      `json:"t"`
		List []string `json:"list"`
	}{
		"ping",
		"identify",
		123,
		[]string{"abc", "def"},
	}

	out := "d2:id8:identify4:listl3:abc3:defe1:q4:ping1:ti123ee"
	if str, err := EncodeStruct(in); err == nil {
		fmt.Println(str, out)
		if str != out {
			t.Fatal(str)
		}
	} else {
		t.Fatal(err)
	}
}

func TestEncodeTop(t *testing.T) {
	in1 := []struct {
		Q    string   `json:"q"`
		Id   string   `json:"id"`
		T    int      `json:"t"`
		List []string `json:"list"`
	}{
		{
			"ping",
			"identify",
			123,
			[]string{"abc", "def"}},
		{
			"r",
			"who",
			321,
			[]string{"rst", "xyz"}},
	}
	out1 := "ld2:id8:identify4:listl3:abc3:defe1:q4:ping1:ti123eed2:id3:who4:listl3:rst3:xyze1:q1:r1:ti321eee"

	cases := []struct {
		in  interface{}
		out string
	}{
		{42, "i42e"},
		{"hello,中国", "12:hello,中国"},
		{[]interface{}{1, "spam", []int{1, 2}}, "li1e4:spamli1ei2eee"},
		{in1, out1},
	}

	for idx := 0; idx != len(cases); idx++ {
		if str, err := Encode(cases[idx].in); err == nil {
			fmt.Println(str, cases[idx].out)
			if str != cases[idx].out {
				t.Fatal(idx, str)
			}
		} else {
			t.Fatal(idx, err)
		}
	}
}

func TestDecodeInt(t *testing.T) {
	in := "i42e"
	var out uint = 42
	var tmp uint = 0

	if err := decodeInteger([]byte(in), reflect.ValueOf(&tmp)); err == nil {
		fmt.Println(tmp, out)
		if tmp != out {
			t.Fatal(tmp, out)
		}
	} else {
		t.Fatal(err)
	}

	in = "i-12345678e"
	var out1 int = -12345678
	var tmp1 int = 0
	if err := decodeInteger([]byte(in), reflect.ValueOf(&tmp1)); err == nil {
		fmt.Println(tmp1, out1)
		if tmp1 != out1 {
			t.Fatal(tmp1, out1)
		}
	} else {
		t.Fatal(err)
	}
}

func TestParseInt(t *testing.T) {
	in := "012345i42eabcdef"
	start := 6
	end := 9
	var out uint32 = 42
	var rettmp uint32 = 0

	retend, err := findInt([]byte(in), start)
	if err != nil {
		t.Fatal(err)
	}
	if retend != end {
		t.Fatal(retend, end)
	}

	retend, err = parseInt([]byte(in), start, reflect.ValueOf(&out))
	if err != nil {
		t.Fatal(err)
	}

	if retend != end && rettmp != out {
		t.Fatal(retend, rettmp)
	}

}

func TestDecodeString(t *testing.T) {
	in := "0123454:spamabcde"
	start := 6
	end := 11
	expect := "spam"
	retstr := ""

	retend, err := findString([]byte(in), start)
	if err != nil {
		t.Fatal(err)
	}
	if retend != end {
		t.Fatal(retend, end)
	}

	retend, err = parseString([]byte(in), start, reflect.ValueOf(&retstr))
	if err != nil {
		t.Fatal(err)
	}
	if retend != end {
		t.Fatal(retend, end)
	}
	if retstr != expect {
		t.Fatal(retstr, in)
	}

	in = "12:hello,中国"
	expect = "hello,中国"
	err = decodeString([]byte(in), reflect.ValueOf(&retstr))
	if err != nil {
		t.Fatal(err)
	}

	if retstr != expect {
		t.Fatal(retstr, in)
	}

}
