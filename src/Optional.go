package ece

import "errors"

type Optional[T any] struct {
	value T
	exist bool
}

func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, exist: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{exist: false}
}

func (o *Optional[T]) Get() (T, error) {
	if !o.exist {
		return o.value, errors.New("value is empty")
	}
	return o.value, nil
}
