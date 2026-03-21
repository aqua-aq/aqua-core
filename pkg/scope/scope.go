package scope

type Scope[K comparable, V any] struct {
	current *scopeFrame[K, V]
}

type scopeFrame[K comparable, V any] struct {
	parent   *scopeFrame[K, V]
	values   map[K]V
	pointers map[K]*scopeFrame[K, V]
}

func New[K comparable, V any]() Scope[K, V] {
	return Scope[K, V]{
		current: &scopeFrame[K, V]{
			values:   make(map[K]V),
			pointers: make(map[K]*scopeFrame[K, V]),
		},
	}
}

func (s Scope[K, V]) Get(key K) (V, bool) {
	if s.current == nil {
		var zero V
		return zero, false
	}
	if frame, ok := s.current.pointers[key]; ok {
		return frame.values[key], true
	}
	var zero V
	return zero, false
}

func (s Scope[K, V]) Has(key K) bool {
	if s.current == nil {
		return false
	}
	_, ok := s.current.pointers[key]
	return ok
}

func (s Scope[K, V]) Push() Scope[K, V] {
	newFrame := &scopeFrame[K, V]{
		parent:   s.current,
		values:   make(map[K]V),
		pointers: make(map[K]*scopeFrame[K, V]),
	}

	for k, v := range s.current.pointers {
		newFrame.pointers[k] = v
	}
	return Scope[K, V]{
		current: newFrame,
	}
}

func (s Scope[K, V]) Pop() (Scope[K, V], bool) {
	if s.current == nil || s.current.parent == nil {
		return Scope[K, V]{nil}, false
	}
	return Scope[K, V]{s.current.parent}, true
}
func (s Scope[K, V]) Rebase() Scope[K, V] {
	s, ok := s.Pop()
	if !ok {
		return New[K, V]()
	}
	return s.Push()
}

func (s Scope[K, V]) Set(key K, val V) {
	s.current.values[key] = val
	s.current.pointers[key] = s.current
}

func (s Scope[K, V]) Delete(key K) bool {
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

func (s Scope[K, V]) All() map[K]V {
	if s.current == nil {
		return nil
	}

	result := make(map[K]V, len(s.current.pointers))

	for k, frame := range s.current.pointers {
		result[k] = frame.values[k]
	}

	return result
}

func (s Scope[K, V]) Clear() {
	clear(s.current.values)
	clear(s.current.pointers)
}
