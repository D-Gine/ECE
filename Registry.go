package ece

import (
	"errors"
	"reflect"
)

type Registry struct {
	components         map[string]any
	events_subscribers map[string][]eventSubscriber
	eventQueue         *EventQueue
}

func NewRegistry() *Registry {
	reg := &Registry{}
	reg.components = make(map[string]any)
	reg.events_subscribers = make(map[string][]eventSubscriber)
	reg.eventQueue = NewEventQueue()
	return reg
}

func typeName[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}

func (r *Registry) RegisterComponent[T any]() *SparseArray[T] {
	return r.RegisterComponentByName[T](typeName[T]())
}

func (r *Registry) RegisterComponentByName[T any](name string) *SparseArray[T] {
	if s, ok := r.components[name]; ok {
		return s.(*SparseArray[T])
	}
	empty := NewSparseArray[T]()
	r.components[name] = empty
	return empty
}

func (r *Registry) GetComponents[T any]() *SparseArray[T] {
	return r.GetComponentsByName[T](typeName[T]())
}

func (r *Registry) GetComponentsByName[T any](name string) *SparseArray[T] {
	s, ok := r.components[name]
	if !ok {
		return nil
	}
	return s.(*SparseArray[T])
}

func (r *Registry) AddComponent[T any](e int, c T) error {
	return r.AddComponentByName(e, typeName[T](), c)
}

func (r *Registry) AddComponentByName[T any](e int, name string, c T) error {
	raw, ok := r.components[name]
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
	r.EmitEventByName(eventName(event), event)
}

func (r *Registry) EmitEventByName(name string, event Event) {
	for _, subscriber := range r.events_subscribers[name] {
		subscriber.callback(r, event)
	}
}

func (r *Registry) EnqueueEvent(event Event) {
	r.eventQueue.Enqueue(event)
}

func (r *Registry) EnqueueEventByName(name string, event Event) {
	r.eventQueue.EnqueueByName(name, event)
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
