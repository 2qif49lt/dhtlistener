package dhtlistener

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Encode(data interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	t := v.Type()
	k := t.Kind()

	switch {
	case reflect.Invalid < k && k <= reflect.Int64:
		return EncodeInt(int(v.Int()))
	case reflect.Uint <= k && k <= reflect.Uint64:
		return EncodeInt(int(v.Uint()))
	case k == reflect.String:
		return EncodeString(v.String())
	case k == reflect.Slice || k == reflect.Array:
		return encodeSlice(v)
	case k == reflect.Map:
		return encodeMap(v)
	case k == reflect.Struct:
		return encodeStruct(v)
	default:
		return "", errors.New("data type no support")
	}
	return "", errors.New("encode fail")
}
func EncodeInt(data int) (string, error) {
	return fmt.Sprintf("i%de", data), nil
}
func EncodeString(data string) (string, error) {
	return strings.Join([]string{strconv.Itoa(len(data)), data}, ":"), nil
}
func encodeSlice(v reflect.Value) (string, error) {
	ret := "l"

	for idx := 0; idx != v.Len(); idx++ {
		itemRet, err := Encode(v.Index(idx).Interface())
		if err != nil {
			return itemRet, err
		}
		ret += itemRet
	}
	ret += "e"

	return ret, nil
}
func EncodeSlice(data interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	return encodeSlice(v)
}

type keySli []reflect.Value

func (sli keySli) Len() int {
	return len(sli)
}
func (sli keySli) Less(i, j int) bool {
	return sli[i].String() < sli[j].String()
}
func (sli keySli) Swap(i, j int) {
	sli[i], sli[j] = sli[j], sli[i]
}

func encodeMap(v reflect.Value) (string, error) {
	ret := "d"
	vkey := v.MapKeys()
	sort.Sort(keySli(vkey))

	for _, val := range vkey {
		itemRet, err := Encode(v.MapIndex(val).Interface())
		if err != nil {
			return itemRet, err
		}
		keyRet, err := EncodeString(val.String())
		if err != nil {
			return keyRet, err
		}
		ret += keyRet + itemRet
	}
	return ret + "e", nil
}
func EncodeMap(data interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	return encodeMap(v)
}
func encodeStruct(v reflect.Value) (string, error) {
	m := make(map[string]interface{})
	t := v.Type()

	for idx := 0; idx != t.NumField(); idx++ {
		if v.Field(idx).CanInterface() {
			name := t.Field(idx).Name
			if tagName := t.Field(idx).Tag.Get("json"); tagName != "" {
				name = tagName
			}

			m[name] = v.Field(idx).Interface()
		}
	}
	return EncodeMap(m)
}

func EncodeStruct(data interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	return encodeStruct(v)
}

// decodeInt decode data string which format is assumed as i...e
func decodeInt(data []byte, v reflect.Value) error {
	numStr := data[1 : len(data)-1]
	num, err := strconv.Atoi(string(numStr))

	v.SetInt(int64(num))
	return err
}

func decodeUint(data []byte, v reflect.Value) error {
	numStr := data[1 : len(data)-1]
	num, err := strconv.Atoi(string(numStr))

	v.SetUint(uint64(num))
	return err
}

func decodeInteger(data []byte, v reflect.Value) error {
	if v.Type().Kind() <= reflect.Int64 {
		return decodeInt(data, v)
	} else {
		return decodeUint(data, v)
	}
}

// decodeString decode data string which format is assumed as num:content
func decodeString(data []byte, v reflect.Value) error {
	idx := strings.Index(string(data), ":")
	if idx == -1 {
		return fmt.Errorf("format wrong")
	}
	nums, err := strconv.Atoi(data[:idx])
	if err != nil {
		return err
	}
	if nums <= 0 {
		return fmt.Errorf("number is invaild")
	}

	if nums != len(data[idx+1:]) {
		return fmt.Errorf("number is not right")
	}
	v.SetString(string(data[idx+1:]))
	return nil
}

func parseInt(data []byte, start int, v reflect.Value) (end int, err error) {
	end = strings.Index(string(data[start+1:]), "e")
	if end == -1 {
		err = errors.New("parseInt donot find end tag")
		return
	}
	end += start + 1 + 1
	err = decodeInteger(data[start:end], v)
	//	num, err = strconv.Atoi(string(data[start+1 : end]))
	return
}
func parseString(data []byte, start int, v reflect.Value) (end int, err error) {
	mid := strings.Index(string(data[start:]), ":")
	if end == -1 {
		err = errors.New("parseInt donot find end tag")
		return
	}
	mid += start

	num := 0
	num, err = strconv.Atoi(string(data[start:mid]))
	if err != nil {
		return
	}
	if len(data[mid+1:]) < num {
		err = errors.New("parseString string length is short")
		return
	}
	end = mid + num + 1
	v.SetString(string(data[mid+1 : end]))
	return
}

// DecodeSlice deco date string which format is assumed as l...e,out'elements are all original byte
func decodeSlice(data []byte, v reflect.Value) error {
	content := data[1 : len(data)-1]

	size := v.Len()
	cur := 0
	s := NewStack()
	idx := 0
	for idx <= len(content) {
		end := 0
		var err error = nil
		switch content[idx] {
		case 'i':
			end, err = parseInt(data, idx, v.Index(cur))
		case 'l':
			{
			}
		case 'd':
			{
			}
		default:
			end, err = parseString(data, idx, v.Index(cur))
		}
		if err != nil {
			return err
		}
		idx = end
		cur++
	}
}

func Decode(data []byte, out interface{}) error {
	v := reflect.ValueOf(out)
	t := v.Type()
	k := t.Kind()

	if k == reflect.Slice {

	} else if k == reflect.Ptr {
		v = v.Elem()
		k = v.Type().Kind()

		switch {
		case reflect.Invalid < k && k <= reflect.Int64:
			return decodeInt(data, v)
		case reflect.Uint <= k && k <= reflect.Uint64:
			return decodeUInt(data, v)
		case k == reflect.String:
			return decodeString(data, v)
		case k == reflect.Slice:
			return decodeSlice(data, v)
		case k == reflect.Map:
			return encodeMap(v)
		case k == reflect.Struct:
			return encodeStruct(v)
		default:
			return "", errors.New("data type no support")
		}
	} else {
		return fmt.Errorf("out type no support")
	}
}
