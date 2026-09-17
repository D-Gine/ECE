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

func (r *Registry) RegisterComponents[T any]() *SparseArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if s, ok := r.components[t]; ok {
		return s.(*SparseArray[T])
	}
	empty := NewSparseArray[T]()
	r.components[t] = empty
	return empty
}

func (r *Registry) GetComponents[T any]() *SparseArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	s, ok := r.components[t]
	if !ok {
		return nil
	}
	return s.(*SparseArray[T])
}

func (r *Registry) AddComponent[T any](e int, c T) error {
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

func (r *Registry) EmitEvent(event Event) {
	t := reflect.TypeOf(event)
	subscribers, ok := r.events_subscribers[t]
	if !ok {
		return
	}
	for _, subscriber := range subscribers {
		subscriber.callback(r, event)
	}
}

func (r *Registry) EnqueueEvent(event Event) {
	r.eventQueue.Enqueue(event)
}

func (r *Registry) ProcessManyEvents(n int) {
	r.eventQueue.ProcessManyEvents(r, n)
}

func (r *Registry) ProcessAllEvents() {
	r.eventQueue.ProcessAllEvents(r)
}

func (r *Registry) ProcessEvents() {
	r.eventQueue.Process(r)
}
