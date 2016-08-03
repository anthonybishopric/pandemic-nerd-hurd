package pandemic

import (
	"errors"
)

type Stack struct {
	vs []interface{}
}

func NewStack() *Stack {
	return &Stack{vs: make([]interface{}, 0)}
}

func (s *Stack) Push(v interface{}) *Stack {
	s.vs = append([]interface{}{v}, s.vs...)
	return s
}

func (s *Stack) Pop() (interface{}, error) {
	l := len(s.vs)
	if l == 0 {
		return nil, errors.New("Empty stack")
	}
	head := s.vs[0]
	s.vs = s.vs[1:]
	return head, nil
}

func (s *Stack) Peek() interface{} {
	l := len(s.vs)
	if l == 0 {
		return nil
	}
	return s.vs[0]
}
