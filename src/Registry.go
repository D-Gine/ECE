package ece

import (
	"errors"
	"reflect"
)

type Registry struct {
	components         map[reflect.Type]any
	events_subscribers map[reflect.Type][]eventSubscriber
	eventQueue         *EventQueue
}

func NewRegistry() *Registry {
	reg := &Registry{}
	reg.components = make(map[reflect.Type]any)
	reg.events_subscribers = make(map[reflect.Type][]eventSubscriber)
	reg.eventQueue = NewEventQueue()
	//	reg.events = make(map[reflect.Type][]func(*Registry, any))
	return reg
}

func RegisterComponents[T any](r *Registry) *SparseArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if s, ok := r.components[t]; ok {
		return s.(*SparseArray[T])
	}
	empty := NewSparseArray[T]()
	r.components[t] = empty
	return empty
}

func GetComponents[T any](r *Registry) *SparseArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	s, ok := r.components[t]
	if !ok {
		return nil
	}
	return s.(*SparseArray[T])
}

func AddComponent[T any](r *Registry, e int, c T) error {
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
