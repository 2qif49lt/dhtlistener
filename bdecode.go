package dhtlistener

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func getStr(data []byte) (str string, err error) {
	mid := strings.Index(string(data[:]), ":")
	if mid == -1 {
		err = errors.New("getStr donot find tag :")
		return
	}

	num := 0
	num, err = strconv.Atoi(string(data[:mid]))
	if err != nil {
		return
	}
	if len(data[mid+1:]) < num {
		err = errors.New("getStr string length is short")
		return
	}
	str = string(data[mid+1:])

	return
}

func getInt(data []byte) (num int, err error) {
	if data[len(data)-1] != 'e' {
		err = errors.New("getInt donot find tag e")
	}
	numStr := data[1 : len(data)-1]
	num, err = strconv.Atoi(string(numStr))

	return
}
func decodex(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("type dont support expect ptr, now:%s", rv.Kind())
	}
	rv = rv.Elem()
	rk := rv.Kind()

	switch data[0] {
	case 'i':
		num, err := getInt(data)
		if err != nil {
			return err
		}
		if rk == reflect.Interface {
			rv.Set(reflect.ValueOf(num))
		} else if rk <= reflect.Int64 {
			rv.SetInt(int64(num))
		} else {
			rv.SetUint(uint64(num))
		}
	case 'l':
		if rk == reflect.Interface {
			var x []interface{}
			defer func(p reflect.Value) { p.Set(rv) }(rv)
			rv = reflect.ValueOf(&x).Elem()
		}

		if rk != reflect.Slice {
			return errors.New("need slice")
		}
		ek := rv.Type().Elem() //element type

		idx := 1

		for idx < len(data)-1 {
			_, end, err := findFirstNode(data, idx)
			if err != nil {
				return err
			}
			elemVal := reflect.New(ek)
			elemit := elemVal.Interface()
			err = decodex(data[idx:end+1], elemit)
			if err != nil {
				return err
			}

			rv.Set(reflect.Append(rv, elemVal.Elem()))

			idx = end + 1
		}

	case 'd':
		if rv.Kind() == reflect.Interface {
			var x map[string]interface{}
			defer func(p reflect.Value) { p.Set(rv) }(rv)
			rv = reflect.ValueOf(&x).Elem()
		}

		if rk != reflect.Map || rv.Type().Key().Kind() != reflect.String {
			return errors.New("need map or map key need be string")
		}
		et := rv.Type().Elem() //element type

		idx := 1
		bkey := true
		key := ""
		keyval := reflect.ValueOf(&key).Elem()

		for idx < len(data)-1 {
			typeid, end, err := findFirstNode(data, idx)
			if err != nil {
				return err
			}
			if bkey && typeid != bencode_type_str {
				return errors.New("expect type string")
			}
			if bkey {
				keystr := ""
				if keystr, err = getStr(data[idx : end+1]); err != nil {
					return err
				}
				keyval.SetString(keystr)
			} else {
				elemVal := reflect.New(et)
				elemit := elemVal.Interface()
				err = decodex(data[idx:end+1], elemit)
				if err != nil {
					return err
				}
				rv.SetMapIndex(keyval, elemVal.Elem())
			}

			idx = end + 1
			bkey = !bkey
		}
	default:
		str, err := getStr(data)
		if err != nil {
			return err
		}
		rv.SetString(str)
	}

	return nil
}

func findFirstNode(data []byte, start int) (typeid, end int, err error) {
	idx := start
	size := len(data)

	s := NewStack()

	for idx != size {
		switch data[idx] {
		case 'i':
			typeid = bencode_type_num
			end, err = findInt(data, idx)
		case 'l':
			s.Push(&typeIdx{bencode_type_list, start, -1})
			_, end, err = findFirstNode(data, idx+1)
		case 'd':
			s.Push(&typeIdx{bencode_type_map, start, -1})
			_, end, err = findFirstNode(data, idx+1)
		case 'e':
			if sitem := s.Pop(); sitem != nil {
				end = idx
				typeid = sitem.typeid
			}
		default:
			typeid = bencode_type_str
			end, err = findString(data, idx)

		}
		if err != nil {
			return
		}
		if s.Size() == 0 {
			return
		}
		idx = end + 1

	}
	return 0, end, errors.New("findType can not find end tag")
}

func findString(data []byte, start int) (end int, err error) {
	mid := strings.Index(string(data[start:]), ":")
	if end == -1 {
		err = errors.New("findString donot find tag :")
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

	return
}

func findInt(data []byte, start int) (end int, err error) {
	end = strings.Index(string(data[start:]), "e")
	if end == -1 {
		err = errors.New("findInt donot find end tag")
		return
	}
	end += start

	return
}
