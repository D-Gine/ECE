package ece

type SparseArray[T any] struct {
	dense  []T
	sparse []Optional[uint64]
}

func (s *SparseArray[T]) Insert(e uint64, c T) {
	back_id := len(s.dense)
	s.sparse[e] = Some(uint64(back_id))
	s.dense = append(s.dense, c)
}

func (s *SparseArray[T]) Get(e uint64) (T, error) {
	idx, err := s.sparse[e].Get()
	if err != nil {
		var null T
		return null, err
	}
	return s.dense[idx], nil
}
