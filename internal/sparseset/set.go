package sparseset

// Store is a type-erased interface implemented by [Set], allowing heterogeneous
// component sets (e.g., Set[Position], Set[Velocity]) to be stored and managed
// together without knowing their concrete component type T.
type Store interface {
	Has(e Entity) bool
	Remove(e Entity)
	Clear()
	Len() int
}

// Set is the fundamental data type. It is backed by three slices,
// sparse, denseID, and dense. See the package documentation for a full
// explanation and examples.
type Set[T any] struct {
	sparse  []int
	denseID []Entity
	dense   []T
}

// New creates a Sparse Set for holding components of type T. A sparse set
// has three underlying slices, sparse, denseID, and dense. The first slice,
// sparse, has length and capacity of initialCap. The other two, denseID and
// dense, have a capacity of initialCap and a length of 0.
// It will panic if initialCap < 0.
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

// Has returns true if the entity specified by the entity handle e
// exists in this set and false otherwise.
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

// Get returns a pointer to the component in this set associated with the
// entity handle e, or nil if e is not present. The returned pointer points
// directly into the backing slice and should not be retained across calls to
// [Set.Insert], [Set.Remove], or [Set.Clear].
func (s *Set[T]) Get(e Entity) *T {
	if !s.Has(e) {
		return nil
	}
	return &s.dense[s.sparse[e.Index()]]
}

// Insert places a copy of val into this set, associated with handle e,
// overwriting any existing component for e and growing the set's capacity if
// needed. Overwriting an existing component preserves its location in dense.
// If Insert() is used to update an existing [Entity], passing a stale
// [Entity] handle will overwrite the sparse slot for that index, corrupting
// the mapping for any live entity sharing that index.
// It will panic if e is [NullEntity].
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

// Remove deletes the component associated with the entity handle e from
// this set. It is safe (and a no-op) to call Remove() with a stale or
// invalid [Entity]. Removal is handled via swap-and-pop, which gives O(1)
// removal, but does not preserve order for [Set.Entities] and [Set.Values].
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

// Clear removes all components from this set. It does retain the allocated
// memory to prevent the need to garbage collect the slice or to reallocate
// the slice in the future. Any components it contains are nulled out so that
// they can be garbage collected.
func (s *Set[T]) Clear() {
	for _, e := range s.denseID {
		s.sparse[e.Index()] = -1
	}
	clear(s.dense)
	s.dense = s.dense[:0]
	s.denseID = s.denseID[:0]
}

// Len returns the number of components contained in this set.
func (s *Set[T]) Len() int {
	return len(s.dense)
}

// Entities returns a contiguous backing slice of the entity handles that have
// a component in this set. It is aligned 1:1 with [Set.Values]. Callers must
// not modify the returned slice, and it must not be retained across calls
// to [Set.Insert], [Set.Remove], or [Set.Clear].
func (s *Set[T]) Entities() []Entity {
	return s.denseID
}

// Values returns the contiguous backing slice of components in this set,
// aligned 1:1 with [Set.Entities]. Elements may be mutated in place, but
// callers must not append to or reslice the returned slice, nor retain it
// across calls to [Set.Insert], [Set.Remove], or [Set.Clear].
func (s *Set[T]) Values() []T {
	return s.dense
}
