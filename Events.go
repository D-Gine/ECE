package ece

import (
	"reflect"
	"sort"
)

type Event interface{}

type eventSubscriber struct {
	callback func(*Registry, Event) any
	priority int
}

func eventName(event Event) string {
	t := reflect.TypeOf(event)
	if t == nil { // interface nil
		return ""
	}
	return t.String()
}

type BreakpointEvent struct{}

func (r *Registry) SubscribeToEvent[T Event](callback func(*Registry, T) any, priority int) {
	r.SubscribeToEventByName(typeName[T](), func(r *Registry, event Event) any {
		return callback(r, event.(T))
	}, priority)
}

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

type CoucouEvent struct {
	text string
}
