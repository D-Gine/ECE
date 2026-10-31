package ece

import (
	"reflect"
	"sort"
)

// Event is the type of every message that flows through the registry.
// It is an empty interface, so any Go value can be used as an event;
// the usual choice is a small struct (or a pointer to one) that carries
// the event's data:
//
//	type DamageEvent struct {
//		Target int
//		Amount int
//	}
//
// Events are routed to subscribers by name (see eventName). The name
// of a value type and of its pointer type differ ("pkg.DamageEvent" vs
// "*pkg.DamageEvent"), so always emit and subscribe using the same form.
type Event interface{}

// eventSubscriber is a callback registered for one event name, together
// with the priority that decides when it runs relative to the other
// subscribers of that event.
type eventSubscriber struct {
	callback func(*Registry, Event) any
	priority int
}

// eventName returns the routing name of an event, which is the string
// representation of its dynamic Go type (for example "main.DamageEvent"
// or "*main.DamageEvent"). It returns "" if event is a nil interface.
func eventName(event Event) string {
	t := reflect.TypeOf(event)
	if t == nil {
		return ""
	}
	return t.String()
}

// BreakpointEvent is a special marker event used to pause the processing
// of an EventQueue. It carries no data.
//
// Enqueue it as a pointer (&BreakpointEvent{}): the queue detects
// breakpoints with a *BreakpointEvent type assertion. See
// EventQueue.Process, EventQueue.ProcessAllEvents and
// EventQueue.ProcessManyEvents for how each one handles it.
type BreakpointEvent struct{}

// SubscribeToEvent registers callback to be called every time an event of
// type T is emitted on this registry. The event name is derived from T,
// and the callback receives the event already converted to T, so no type
// assertion is needed in user code.
//
// priority controls the call order among the subscribers of the same
// event: lower values run first, and subscribers with equal priority run
// in the order they were registered.
//
//	r.SubscribeToEvent[DamageEvent](onDamage, 0)
func (r *Registry) SubscribeToEvent[T Event](callback func(*Registry, T) any, priority int) {
	r.SubscribeToEventByName(typeName[T](), func(r *Registry, event Event) any {
		return callback(r, event.(T))
	}, priority)
}

// SubscribeToEventByName registers callback for the event identified by
// name. Use it when events are routed by a custom name (see
// Registry.EmitEventByName) instead of by Go type.
//
// callback receives the event as a plain Event and must type-assert it
// itself. priority has the same meaning as in SubscribeToEvent: lower
// values run first, ties keep registration order.
func (r *Registry) SubscribeToEventByName(name string, callback func(*Registry, Event) any, priority int) {
	wrapped := func(r *Registry, event Event) any {
		return callback(r, event)
	}
	r.events_subscribers[name] = append(r.events_subscribers[name], eventSubscriber{
		callback: wrapped,
		priority: priority,
	})
	sort.SliceStable(r.events_subscribers[name], func(i, j int) bool {
		return r.events_subscribers[name][i].priority < r.events_subscribers[name][j].priority
	})
}
