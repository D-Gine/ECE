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

type BreakpointEvent struct{}

func subscribe[T Event](r *Registry, callback func(*Registry, T) any, priority int) {
	t := reflect.TypeFor[T]()
	wrapped := func(r *Registry, event Event) any {
		return callback(r, event.(T))
	}
	r.events_subscribers[t] = append(r.events_subscribers[t], eventSubscriber{
		callback: wrapped,
		priority: priority,
	})
	sort.SliceStable(r.events_subscribers[t], func(i, j int) bool {
		return r.events_subscribers[t][i].priority < r.events_subscribers[t][j].priority
	})
}

type CoucouEvent struct {
	text string
}

func onCoucouEvent(r *Registry, e *CoucouEvent) any {
	println("CoucouEvent received with text:", e.text)
	return nil
}
