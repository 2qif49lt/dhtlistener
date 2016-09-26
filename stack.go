package dhtlistener

type stacknode struct {
	value interface{}
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

func (s *Stack) Push(v interface{}) {
	n := &stacknode{
		value: v,
		next:  s.top,
	}
	s.top = n
	s.size++
}

func (s *Stack) Pop() interface{} {
	if s.size == 0 {
		return nil
	}
	n := s.top
	s.size--
	s.top = s.top.next
	return n.value
}

func (s *Stack) Peek() interface{} {
	if s.size == 0 {
		return nil
	}

	return s.top.value
}
