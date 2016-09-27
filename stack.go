package dhtlistener

type typeIdx struct {
	typeid int
	beg    int
	end    int
}

type stacknode struct {
	value *typeIdx
	next  *stacknode
}
type Stack struct {
	top  *stacknode
	size int
}

func NewStack() *Stack {
	return &Stack{
		top:  nil,
		size: 0,
	}
}

func (s *Stack) Size() int {
	return s.size
}

func (s *Stack) Push(v *typeIdx) {
	n := &stacknode{
		value: v,
		next:  s.top,
	}
	s.top = n
	s.size++
}

func (s *Stack) Pop() *typeIdx {
	if s.size == 0 {
		return nil
	}
	n := s.top
	s.size--
	s.top = s.top.next
	return n.value
}

func (s *Stack) Peek() *typeIdx {
	if s.size == 0 {
		return nil
	}

	return s.top.value
}

const (
	bencode_type_num = iota
	bencode_type_str
	bencode_type_list
	bencode_type_map
)
