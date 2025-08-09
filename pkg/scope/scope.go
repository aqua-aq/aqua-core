package scope

type Scope[T any] struct {
	current *scopeFrame[T]
}

type scopeFrame[T any] struct {
	parent   *scopeFrame[T]
	values   map[string]T
	pointers map[string]*scopeFrame[T]
}

func New[T any]() Scope[T] {
	return Scope[T]{
		current: &scopeFrame[T]{
			values:   make(map[string]T),
			pointers: make(map[string]*scopeFrame[T]),
		},
	}
}

func (s Scope[T]) Get(key string) (T, bool) {
	if s.current == nil {
		var zero T
		return zero, false
	}
	if frame, ok := s.current.pointers[key]; ok {
		return frame.values[key], true
	}
	var zero T
	return zero, false
}

func (s Scope[T]) Has(key string) bool {
	if s.current == nil {
		return false
	}
	_, ok := s.current.pointers[key]
	return ok
}

func (s Scope[T]) Push() Scope[T] {
	newFrame := &scopeFrame[T]{
		parent:   s.current,
		values:   make(map[string]T),
		pointers: make(map[string]*scopeFrame[T]),
	}

	for k, v := range s.current.pointers {
		newFrame.pointers[k] = v
	}
	return Scope[T]{
		current: newFrame,
	}
}

func (s Scope[T]) Pop() (Scope[T], bool) {
	if s.current == nil || s.current.parent == nil {
		return Scope[T]{nil}, false
	}
	return Scope[T]{s.current.parent}, true
}
func (s Scope[T]) Rebase() Scope[T] {
	s, ok := s.Pop()
	if !ok {
		return New[T]()
	}
	return s.Push()
}

func (s Scope[T]) Set(key string, val T) {
	s.current.values[key] = val
	s.current.pointers[key] = s.current
}

func (s Scope[T]) Delete(key string) bool {
	if s.current == nil {
		return false
	}

	_, hasValue := s.current.values[key]
	_, hasPointer := s.current.pointers[key]

	if hasValue || hasPointer {
		delete(s.current.values, key)
		delete(s.current.pointers, key)
		return true
	}

	return false
}

func (s Scope[T]) Clear() {
	clear(s.current.values)
	clear(s.current.pointers)
}
