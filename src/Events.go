package ece

import "reflect"

type EventFunc func(*Registry, any)

type EventsQueue struct {
	t reflect.Type
	data any
}

type Events struct {
	list  map[reflect.Type][]EventFunc
	queue []EventsQueue
}

func NewEvents() *Events {
	e := &Events{}
	e.list = make(map[reflect.Type][]EventFunc)
	return e
}

func Subscribe[T any](r *Registry, fn func(*Registry, T)) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	wrapped := func(reg *Registry, data any) {
		typed, ok := data.(T)
		if !ok {
			return
		}
		fn(reg, typed)
	}
	r.events.list[t] = append(r.events.list[t], wrapped)
}

func Unsubscribe[T any](r *Registry) {

}

func Fire[T any](r *Registry, data T) {

}

func AddToQueue[T any](r *Registry, data T) {
	r.events.queue = append(
		r.events.queue,
		EventsQueue{t, data})
}

func PlayQueue(r *Registry) {
	for _, obj := range r.events.queue {
		typed, ok := obj.data.()
		if !ok {
			continue
		} else {
			Fire(r, typed)
		}
	}
}
