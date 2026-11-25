package ece

import "errors"

// Optional is a minimal "maybe" type (similar to Option in Rust or
// std::optional in C++). It either holds a value of type T ("Some")
// or holds nothing ("None").
//
// Always build an Optional with Some or None rather than with a struct
// literal, so that the exist flag is always consistent with the value.
//
// The zero value of Optional[T] is equivalent to None[T]().
type Optional[T any] struct {
	value T    // the stored value; only meaningful when exist is true
	exist bool // true if the Optional holds a value
}

// Some returns an Optional that holds the value v.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, exist: true}
}

// None returns an empty Optional of type T.
// The type parameter cannot be inferred, so it must be given explicitly:
//
//	empty := None[int]()
func None[T any]() Optional[T] {
	return Optional[T]{exist: false}
}

// Get returns the stored value and a nil error if the Optional holds a
// value. If the Optional is empty, it returns the zero value of T and
// the error "value is empty".
//
// Always check the returned error before using the value.
func (o *Optional[T]) Get() (T, error) {
	if !o.exist {
		return o.value, errors.New("value is empty")
	}
	return o.value, nil
}
