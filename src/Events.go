package ece

import (
	"reflect"
	"sort"
)

type Event interface {
	triggerEvent(*Registry)
}

type eventSubscriber struct {
	callback func(*Registry, Event)
	priority int
}

func subscribe[T Event](r *Registry, callback func(*Registry, T), priority int) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	wrapped := func(r *Registry, event Event) {
		callback(r, event.(T))
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

func (e *CoucouEvent) triggerEvent(r *Registry) {
	t := reflect.TypeOf(e)
	subscribers, ok := r.events_subscribers[t]
	if !ok {
		return
	}
	for _, subscriber := range subscribers {
		subscriber.callback(r, e)
	}
}

func onCoucouEvent(r *Registry, e *CoucouEvent) {
	println("CoucouEvent received with text:", e.text)
}
