package stack

// Element is an element in a single linked list
type Element struct {
	Value interface{}
	next  *Element
}

// Next accesses the next element connected to this current node
func (e *Element) Next() *Element {
	return e.next
}

// Stack represents a stack
// The zero value for Stack is an unintialized stack
// ought to initialize by invoking Init()
type Stack struct {
	top  *Element // top element of the stack
	size int      // size of the stack
}

// New returns an initialized stack
func New() *Stack {
	return &Stack{nil, 0}
}

// Push adds an element to the stack
func (s *Stack) Push(x interface{}) {
	s.top = &Element{x, s.top}
	s.size++
}

// Pop out the top element from the stack
func (s *Stack) Pop() interface{} {
	var x interface{}
	if nil != s.top {
		x, s.top = s.top.Value, s.top.next
		s.size--
	}

	return x
}

// Peek returns the element in the top of the s
func (s *Stack) Peek() interface{} {
	if nil != s.top {
		return s.top.Value
	}

	return nil
}

// Peek returns the two elements in the top of the s
func (s *Stack) Peek2() (interface{}, interface{}) {
	var x, y interface{}

	if nil != s.top {
		x = s.top.Value

		if nil != s.top.next {
			y = s.top.next.Value
		}
	}

	return x, y
}

// Len returns the size of the s
func (s *Stack) Len() int {
	return s.size
}

// Empty returns true if the stack is empty
func (s *Stack) Empty() bool {
	return 0 == s.size
}

// ValueSlice returns all elements of stack in slice
func (s *Stack) ValueSlice() []interface{}{
	vs := make([]interface{}, s.size)
	ele := s.top
	for i := s.size - 1; i >= 0; i--{
		vs[i] = ele.Value
		ele = ele.next
	}
	return vs
}
