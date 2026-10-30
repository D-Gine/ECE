package ece

import (
	"errors"
)

// INIT_ARR_SIZE is the number of slots allocated in the sparse array
// the first time an element is inserted.
const INIT_ARR_SIZE int = 4

// DEFAULT_ARR_SIZE is the sentinel value of SparseArray.size meaning
// "the sparse array has not been allocated yet".
const DEFAULT_ARR_SIZE int = -1

// SparseArray is a sparse set that maps an entity ID (a non-negative int)
// to a component of type T. In ECE, each component type gets its own
// SparseArray[T].
//
// It uses three parallel structures:
//
//	sparse          indexed by entity ID; holds the position of that
//	                entity's component in dense (or None if it has none).
//	dense           tightly packed component values, no holes. Good for
//	                cache-friendly iteration.
//	reverse         reverse mapping: reverse[i] is the entity ID
//	                that owns dense[i]. Needed so Remove can fix up the
//	                sparse entry of the element it moves.
//
// Invariant: for every live entity e, with i := sparse[e],
// dense[i] is its component and reverse[i] == e.
//
// Complexity: Get and Remove are O(1). Insert is amortized O(1) (the
// sparse array doubles when it needs to grow).
//
// SparseArray is NOT safe for concurrent use; guard it with a mutex if
// several goroutines access it.
type SparseArray[T any] struct {
	dense   []T             // packed component values
	reverse []int           // dense index -> entity ID
	sparse  []Optional[int] // entity ID -> dense index
	size    int             // current length of sparse, or DEFAULT_ARR_SIZE if unallocated
}

// NewSparseArray creates an empty SparseArray that stores components of
// type T and returns a pointer to it. No memory is allocated for the
// sparse array until the first Insert.
//
//	positions := ece.NewSparseArray[Vec2]()
func NewSparseArray[T any]() *SparseArray[T] {
	s := &SparseArray[T]{}
	s.size = DEFAULT_ARR_SIZE
	return s
}

// Insert attaches component c to entity e.
//
// If e is larger than the current sparse capacity, the sparse array is
// grown (starting at INIT_ARR_SIZE, then doubling) until it can hold e.
// The new component is appended at the end of the dense array.
//
// Caveats for callers:
//   - e must be >= 0, otherwise Insert panics (negative slice index).
//   - Insert does not check whether e already has a component. Inserting
//     twice for the same entity leaves an orphaned value in the dense
//     array and breaks the invariant. Call Remove(e) first if you want to
//     replace an existing component.
//   - Memory usage of the sparse array is proportional to the largest
//     entity ID, so prefer small, recycled IDs.
func (s *SparseArray[T]) Insert(e int, c T) {
	back_id := len(s.dense)

	if s.size-1 < e {
		for s.size-1 < e {
			if s.size == DEFAULT_ARR_SIZE {
				for i := 0; i < INIT_ARR_SIZE; i++ {
					s.sparse = append(s.sparse, None[int]())
				}
				s.size = INIT_ARR_SIZE
			} else {
				upgrade := make([]Optional[int], s.size)
				for i := range s.size {
					upgrade[i] = None[int]()
				}
				s.sparse = append(s.sparse, upgrade...)
				s.size += len(upgrade)
			}
		}
	}

	s.sparse[e] = Some[int](back_id)
	s.dense = append(s.dense, c)
	s.reverse = append(s.reverse, e)
}

// Get returns the component attached to entity e.
//
// If e has no component, it returns the zero value of T and the error
// "value is empty".
//
// Caveat: Get does not bounds-check e. Calling it with a negative ID, or
// an ID beyond the sparse array (never inserted, nor below a previously
// inserted ID), panics with an index out of range.
func (s *SparseArray[T]) Get(e int) (T, error) {
	idx, err := s.sparse[e].Get()
	if err != nil {
		var null T
		return null, err
	}
	return s.dense[idx], nil
}

// Remove detaches the component of entity e, if any.
//
// Removing an entity that is out of range or has no component is a no-op
// and returns nil. An error is returned only if the internal state is
// corrupted (the sparse entry points outside the dense array).
//
// To keep the dense array free of holes, Remove uses "swap and pop": the
// last dense element is moved into the freed slot, its sparse entry is
// updated, and the dense array is shortened by one. As a consequence,
// the iteration order of the dense array is NOT stable across removals.
func (s *SparseArray[T]) Remove(e int) error {
	if e < 0 || e >= len(s.sparse) {
		return nil
	}
	to_del, err := s.sparse[e].Get()
	if err != nil {
		return nil
	}
	if to_del < 0 || to_del >= len(s.dense) {
		return errors.New("error: entity does not exist")
	}

	last := len(s.dense) - 1
	s.sparse[e] = None[int]()

	if to_del != last {
		movedEntity := s.reverse[last]
		s.dense[to_del] = s.dense[last]
		s.reverse[to_del] = movedEntity
		s.sparse[movedEntity] = Some(to_del)
	}

	s.dense = s.dense[:last]
	s.reverse = s.reverse[:last]
	return nil
}
