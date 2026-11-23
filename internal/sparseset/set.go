package sparseset

type Store interface {
	Has(e Entity) bool
	Remove(e Entity)
	Clear()
	Len() int
}

type Set[T any] struct {
	sparse  []int
	denseID []Entity
	dense   []T
}

func New[T any](initialCap int) *Set[T] {
	sparse := make([]int, initialCap)
	for i := range sparse {
		sparse[i] = -1
	}
	return &Set[T]{
		sparse:  sparse,
		denseID: make([]Entity, 0, initialCap),
		dense:   make([]T, 0, initialCap),
	}
}

func (s *Set[T]) Has(e Entity) bool {
	if e == NullEntity {
		return false
	}
	idx := int(e.Index())
	if idx >= len(s.sparse) {
		return false
	}
	denseIdx := s.sparse[idx]
	return denseIdx >= 0 && denseIdx < len(s.dense) && s.denseID[denseIdx] == e
}

func (s *Set[T]) Get(e Entity) *T {
	if !s.Has(e) {
		return nil
	}
	return &s.dense[s.sparse[e.Index()]]
}

func (s *Set[T]) Insert(e Entity, val T) {
	if e == NullEntity {
		panic("sparseset: cannot insert NullEntity")
	}

	idx := int(e.Index())

	if s.Has(e) {
		s.dense[s.sparse[idx]] = val
		return
	}

	if idx >= len(s.sparse) {
		newLen := max(idx+1, len(s.sparse)*2)
		newSparse := make([]int, newLen)
		copy(newSparse, s.sparse)
		for i := len(s.sparse); i < newLen; i++ {
			newSparse[i] = -1
		}
		s.sparse = newSparse
	}

	s.sparse[idx] = len(s.dense)
	s.denseID = append(s.denseID, e)
	s.dense = append(s.dense, val)
}

func (s *Set[T]) Remove(e Entity) {
	if !s.Has(e) {
		return
	}

	idx := int(e.Index())
	denseIdx := s.sparse[idx]
	lastIdx := len(s.dense) - 1
	lastEntity := s.denseID[lastIdx]

	s.dense[denseIdx] = s.dense[lastIdx]
	s.denseID[denseIdx] = lastEntity
	s.sparse[lastEntity.Index()] = denseIdx

	var zero T
	s.dense[lastIdx] = zero

	s.dense = s.dense[:lastIdx]
	s.denseID = s.denseID[:lastIdx]
	s.sparse[idx] = -1
}

func (s *Set[T]) Clear() {
	for _, e := range s.denseID {
		s.sparse[e.Index()] = -1
	}
	clear(s.dense)
	s.dense = s.dense[:0]
	s.denseID = s.denseID[:0]
}

func (s *Set[T]) Len() int {
	return len(s.dense)
}

func (s *Set[T]) Entities() []Entity {
	return s.denseID
}

func (s *Set[T]) Values() []T {
	return s.dense
}
