package stack

import (
	"fmt"
	"iter"
	"reflect"
	"strings"
)

type Stack[T any] struct {
	vals []T
}

func New[T any]() Stack[T] {
	return Stack[T]{vals: []T{}}
}

func Cap[T any](cap int) Stack[T] {
	return Stack[T]{vals: make([]T, 0, cap)}
}

func (s *Stack[T]) Push(vals ...T) {
	if s.vals == nil {
		s.vals = vals
		return
	}
	s.vals = append(s.vals, vals...)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}
	last := len(s.vals) - 1
	val := s.vals[last]
	s.vals = s.vals[:last]
	return val, true
}

func (s *Stack[T]) Cut(n int) []T {
	if n < 0 {
		panic(fmt.Sprintf("Stack[%v].Cut(%v): n must be positive", reflect.TypeFor[T](), n))
	}
	length := len(s.vals)
	if n > length {
		return nil
	}
	mn := len(s.vals) - n
	result := s.vals[mn:]
	s.vals = s.vals[:mn]
	return result
}

func (s *Stack[T]) Swap() bool {
	if len(s.vals) < 2 {
		return false
	}
	last := len(s.vals) - 1
	s.vals[last], s.vals[last-1] = s.vals[last-1], s.vals[last]
	return true
}

func (s *Stack[T]) Dub() bool {
	if len(s.vals) == 0 {
		return false
	}
	s.vals = append(s.vals, s.vals[len(s.vals)-1])
	return true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}
	return s.vals[len(s.vals)-1], true
}

func (s *Stack[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s.vals) - 1; i >= 0; i-- {
			if !yield(s.vals[i]) {
				return
			}
		}
	}
}

func (s *Stack[T]) Clear() {
	s.vals = s.vals[:0]
}

func (s *Stack[T]) Copy(slice []T) {
	copy(slice, s.vals)
}

func (s *Stack[T]) Len() int {
	return len(s.vals)
}

func (s *Stack[T]) String() string {
	sb := strings.Builder{}
	for i := len(s.vals) - 1; i >= 0; i-- {
		sb.WriteString(fmt.Sprintf("%v\n", s.vals[i]))
	}
	return sb.String()
}
