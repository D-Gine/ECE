package ece

import (
	"errors"
	"reflect"
)

// Registry is the central object of ECE. It owns:
//   - the component storages, one SparseArray per component type;
//   - the event subscribers, grouped by event name;
//   - an EventQueue for events whose delivery is deferred.
//
// Entities are plain ints; a component belongs to an entity when it is
// stored in the SparseArray at that entity's ID. There are no systems:
// all logic lives in callbacks subscribed to events.
//
// Always create a Registry with NewRegistry. The zero value is not
// usable because its internal maps are nil.
//
// Registry is NOT safe for concurrent use.
type Registry struct {
	components         map[string]any
	events_subscribers map[string][]eventSubscriber
	eventQueue         *EventQueue
}

// NewRegistry creates an empty Registry, with no components, no
// subscribers and an empty event queue, and returns a pointer to it.
func NewRegistry() *Registry {
	reg := &Registry{}
	reg.components = make(map[string]any)
	reg.events_subscribers = make(map[string][]eventSubscriber)
	reg.eventQueue = NewEventQueue()
	return reg
}

// typeName returns the string representation of type T (for example
// "main.Position"). It is the default key used to store components and
// to route events of that type.
func typeName[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}

// RegisterComponent creates the storage for components of type T and
// returns it. The storage is keyed by the name of T.
//
// It is safe to call more than once: if T is already registered, the
// existing storage is returned and no data is lost.
//
// A component type must be registered before AddComponent is used with it.
func (r *Registry) RegisterComponent[T any]() *SparseArray[T] {
	return r.RegisterComponentByName[T](typeName[T]())
}

// RegisterComponentByName is like RegisterComponent but stores the
// component under a custom name. This allows several independent
// storages of the same Go type, for example two SparseArray[int] named
// "health" and "mana".
//
// If name is already registered, the existing storage is returned. It
// panics if that existing storage was registered with a different type
// parameter.
func (r *Registry) RegisterComponentByName[T any](name string) *SparseArray[T] {
	if s, ok := r.components[name]; ok {
		return s.(*SparseArray[T])
	}
	empty := NewSparseArray[T]()
	r.components[name] = empty
	return empty
}

// GetComponents returns the storage of components of type T, or nil if T
// has not been registered. Always check for nil before using the result.
func (r *Registry) GetComponents[T any]() *SparseArray[T] {
	return r.GetComponentsByName[T](typeName[T]())
}

// GetComponentsByName returns the storage registered under name, or nil
// if there is none. It panics if the storage under that name holds a
// type other than T.
func (r *Registry) GetComponentsByName[T any](name string) *SparseArray[T] {
	s, ok := r.components[name]
	if !ok {
		return nil
	}
	return s.(*SparseArray[T])
}

// AddComponent attaches component c to entity e, in the storage of type
// T. It returns an error if T has not been registered with
// RegisterComponent.
//
// See SparseArray.Insert for the caveats on e (it must be >= 0, and
// adding twice to the same entity is not supported).
func (r *Registry) AddComponent[T any](e int, c T) error {
	return r.AddComponentByName(e, typeName[T](), c)
}

// AddComponentByName attaches component c to entity e, in the storage
// registered under name. It returns an error if name is not registered,
// or if the storage under that name holds a type other than T.
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

// EmitEvent delivers event immediately to every subscriber of its type,
// in priority order. It returns once all callbacks have finished.
//
// The values returned by the callbacks are ignored. Emitting an event
// nobody subscribed to does nothing.
func (r *Registry) EmitEvent(event Event) {
	r.EmitEventByName(eventName(event), event)
}

// EmitEventByName delivers event immediately to every subscriber
// registered under name, in priority order. The values returned by the
// callbacks are ignored.
func (r *Registry) EmitEventByName(name string, event Event) {
	for _, subscriber := range r.events_subscribers[name] {
		subscriber.callback(r, event)
	}
}

// EnqueueEvent adds event to the registry's queue so it is delivered
// later, when one of the Process methods runs. Use EmitEvent instead if
// the event must be handled right away.
func (r *Registry) EnqueueEvent(event Event) {
	r.eventQueue.Enqueue(event)
}

// EnqueueEventByName adds event to the registry's queue under a custom
// routing name. See EmitEventByName.
func (r *Registry) EnqueueEventByName(name string, event Event) {
	r.eventQueue.EnqueueByName(name, event)
}

// ProcessManyEvents emits up to n queued events and returns the result.
// See EventQueue.ProcessManyEvents for the details, including how
// breakpoints are handled.
func (r *Registry) ProcessManyEvents(n int) []Pair[string, Event] {
	return r.eventQueue.ProcessManyEvents(r, n)
}

// ProcessAllEvents emits every queued event, ignoring breakpoints, and
// returns the events that were handled. See EventQueue.ProcessAllEvents.
func (r *Registry) ProcessAllEvents() []Pair[string, Event] {
	return r.eventQueue.ProcessAllEvents(r)
}

// ProcessEvents emits queued events until the queue is empty or a
// *BreakpointEvent is reached, and returns the events that were
// handled. Typical use is one call per frame. See EventQueue.Process.
func (r *Registry) ProcessEvents() []Pair[string, Event] {
	return r.eventQueue.Process(r)
}
