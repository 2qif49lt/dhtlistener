package dhtlistener

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
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

func decodeIntInterface(data []byte, v reflect.Value) error {
	numStr := data[1 : len(data)-1]
	num, err := strconv.Atoi(string(numStr))

	v.Set(reflect.ValueOf(num))
	return err
}

func decodeInteger(data []byte, v reflect.Value) error {
	/*
		if v.Type().Kind() != reflect.Ptr {
			return fmt.Errorf("decodeInteger need point but get %s", v.Type().Kind())
		} else {
			v = v.Elem()
		}
	*/
	if v.Type().Kind() == reflect.Interface {
		return decodeIntInterface(data, v)
	}
	if v.Type().Kind() <= reflect.Int64 {
		return decodeInt(data, v)
	} else {
		return decodeUint(data, v)
	}
}

// parseInt try to obtain number from the begin position start,return end position.
func parseInt(data []byte, start int, v reflect.Value) (end int, err error) {
	fmt.Println(string(data), start)
	end = strings.Index(string(data[start:]), "e")
	if end == -1 {
		err = errors.New("parseInt donot find end tag")
		return
	}
	end += start
	err = decodeInteger(data[start:end+1], v)

	return
}

// decodeString decode data string which format is assumed as num:content
func decodeString(data []byte, v reflect.Value) error {
	idx := strings.Index(string(data), ":")
	if idx == -1 {
		return fmt.Errorf("format wrong")
	}
	nums, err := strconv.Atoi(string(data[:idx]))
	if err != nil {
		return err
	}
	if nums <= 0 {
		return fmt.Errorf("number is invaild")
	}

	if nums != len(data[idx+1:]) {
		return fmt.Errorf("number is not right")
	}
	if v.Kind() == reflect.Interface {
		v.Set(reflect.ValueOf(string(data[idx+1:])))
	} else {
		if v.Kind() != reflect.Ptr {
			return errors.New("decodeString need point")
		} else {
			v = v.Elem()
		}
		v.SetString(string(data[idx+1:]))
	}

	return nil
}

func parseString(data []byte, start int, v reflect.Value) (end int, err error) {
	if v.Type().Kind() != reflect.Ptr {
		err = errors.New("parseString need point")
		return
	} else {
		v = v.Elem()
	}
	mid := strings.Index(string(data[start:]), ":")
	if end == -1 {
		err = errors.New("parseString donot find end tag")
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
	end = mid + num
	v.SetString(string(data[mid+1 : end+1]))
	return
}

func findMap(data []byte, start int) (end int, err error) {
	if data[start] != 'd' {
		err = errors.New("findMap format expect l")
		return
	}
	return findType(bencode_type_map, data, start)

}

// findType fisrt type
func findType(typeid int, data []byte, start int) (end int, err error) {
	s := NewStack()
	s.Push(&typeIdx{typeid, start, -1})

	idx := start + 1
	for idx != len(data[start:]) {
		switch data[idx] {
		case 'i':
			fmt.Println("findType case i")
			end, err = findInt(data, idx)
		case 'l':
			fmt.Println("findType case l")

			s.Push(&typeIdx{bencode_type_list, idx, -1})
			end, err = findSlice(data, idx)
		case 'd':
			//	s.Push(&typeIdx{bencode_type_map, idx, -1})
			end, err = findMap(data, idx)
		case 'e':
			fmt.Println("case e s.Size ", s.Size())
			if sitem := s.Pop(); sitem != nil && s.Size() == 0 {
				end = idx
				fmt.Println("e return")
				return
			}
			fmt.Println("s.Size ", s.Size())

			idx++
			continue
		default:
			end, err = findString(data, idx)
			if err == nil {
				fmt.Println("findType case string")

			}
		}
		fmt.Println(string(data[idx : end+1]))

		if err != nil {
			return
		}
		idx = end + 1

	}
	return end, errors.New("findType can not find end tag")
}
func findSlice(data []byte, start int) (end int, err error) {
	if data[start] != 'l' {
		err = errors.New("findSlice format expect l")
		return
	}

	return findType(bencode_type_list, data, start)
}

// DecodeSlice deco date string which format is assumed as l...e,out'elements are all original byte
func decodeSlice(data []byte, v reflect.Value) error {
	fmt.Println("decodeSlice", string(data), v.Interface())

	size := v.Len()
	cur := 0
	idx := 0
	for cur < size && idx < len(data) {
		typeid, end, err := findFirstNode(data, idx)
		if err != nil {
			return err
		}
		switch typeid {
		case bencode_type_num:
			err = decodeInteger(data[idx:end+1], v.Index(cur))
		case bencode_type_str:
			err = decodeString(data[idx:end+1], v.Index(cur))
		case bencode_type_list:
			err = decodeSlice(data[idx:end+1], v.Index(cur))
		case bencode_type_map:
			err = decodeMap(data[idx:end+1], v.Index(cur))
		default:
			err = fmt.Errorf("type %d dont support", typeid)
		}
		if err != nil {
			return err
		}
		idx = end + 1
	}
	return nil
}

func decodeMap(data []byte, v reflect.Value) error {
	if v.Type().Key().Kind() != reflect.String && v.Type().Key().Kind() != reflect.Interface {
		return fmt.Errorf("decodeMap expect string or interface ")
	}
	return nil
}

func decodeType(data []byte, v reflect.Value) error {
	if v.Type().Kind() != reflect.Ptr {
		return errors.New("decodeType need point")
	} else {
		v = v.Elem()
	}

	size := v.Len()
	cur := 0
	idx := 0

	for idx <= len(data) && cur < size {
		end := 0
		var err error = nil
		switch data[idx] {
		case 'i':
			fmt.Println("decodeType case i", idx)
			end, err = parseInt(data, idx, v.Index(cur))
		case 'l':
			if end, err = findSlice(data, idx); err == nil {
				err = decodeSlice(data[idx:end+1], v.Index(cur))
			}
		case 'd':
			if end, err = findMap(data, idx); err == nil {
				err = decodeMap(data, v.Index(cur))
			}
		default:
			end, err = parseString(data, idx, v.Index(cur))
		}

		if err != nil {
			return err
		}

		idx = end + 1
		cur++
	}
	return nil
}
