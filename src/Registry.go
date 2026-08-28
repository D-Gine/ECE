package ece

import (
	"errors"
	"reflect"
)

type Registry struct {
	components map[reflect.Type]any
}

func NewRegistry() *Registry {
	return &Registry{}
}

func RegisterComponents[T any](r *Registry) *SparseArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if s, ok := r.components[t]; ok {
		return s.(*SparseArray[T])
	}
	s := &SparseArray[T]{}
	r.components[t] = s
	return s
}

func GetComponents[T any](r *Registry) *SparseArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	s, ok := r.components[t]
	if !ok {
		return nil
	}
	return s.(*SparseArray[T])
}

func AddComponent[T any](r *Registry, e uint64, c T) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	raw, ok := r.components[t]
	if !ok {
		return errors.New("not existing type in registry")
	}
	s, ok := raw.(*SparseArray[T])
	if !ok {
		return errors.New("stored component has unexpected type")
	}
	s.Insert(e, c)
	return nil
}
