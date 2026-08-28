package ece

import "errors"

var NoneErr = errors.New("value is empty")

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
		return o.value, NoneErr
	}
	return o.value, nil
}
