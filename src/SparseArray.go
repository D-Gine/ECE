package ece

import (
	"errors"
)

const INIT_ARR_SIZE int = 4
const DEFAULT_ARR_SIZE int = -1

type SparseArray[T any] struct {
	dense           []T
	dense_to_sparse []int
	sparse          []Optional[int]
	size            int
}

func NewSparseArray[T any]() *SparseArray[T] {
	s := &SparseArray[T]{}
	s.size = DEFAULT_ARR_SIZE
	return s
}

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
	s.dense_to_sparse = append(s.dense_to_sparse, e)
}

func (s *SparseArray[T]) Get(e int) (T, error) {
	idx, err := s.sparse[e].Get()
	if err != nil {
		var null T
		return null, err
	}
	return s.dense[idx], nil
}

func (s *SparseArray[T]) Remove(e int) error {
	to_del, err := s.sparse[e].Get()
	if err != nil {
		return nil
	}
	if len(s.dense) < to_del {
		return errors.New("error: entity does not exist")
	}
	last := len(s.dense) - 1
	s.sparse[e] = None[int]()
	s.sparse[s.dense_to_sparse[last]] = Some(to_del)
	s.dense[to_del], s.dense[last] = s.dense[last], s.dense[to_del]
	s.dense = s.dense[:len(s.dense)-1]
	return nil
}
