package ece

// Pair holds two values together. In this package it is used to store an
// event with its routing name: First is the name and Second is the event.
type Pair[T string, U any] struct {
	First  T
	Second U
}

// EventQueue is a FIFO queue of events waiting to be emitted. It lets
// events be recorded now and dispatched later, for example once per
// frame of the game loop.
//
// A queue is normally owned by a Registry and used through the
// Registry.EnqueueEvent and Registry.Process* methods.
//
// EventQueue is NOT safe for concurrent use.
type EventQueue struct {
	events []Pair[string, Event]
}

// NewEventQueue creates an empty EventQueue and returns a pointer to it.
func NewEventQueue() *EventQueue {
	return &EventQueue{
		events: make([]Pair[string, Event], 0),
	}
}

// Enqueue adds event at the back of the queue. Its routing name is
// derived from its Go type (see eventName).
func (eq *EventQueue) Enqueue(event Event) {
	eq.EnqueueByName(eventName(event), event)
}

// EnqueueByName adds event at the back of the queue under a custom
// routing name. When the event is emitted, it is delivered to the
// subscribers registered for name.
func (eq *EventQueue) EnqueueByName(name string, event Event) {
	eq.events = append(eq.events, Pair[string, Event]{name, event})
}

// Dequeue removes the event at the front of the queue and returns it
// with its name. If the queue is empty, it returns a zero Pair (empty
// name and nil event); use IsEmpty to tell the two cases apart.
func (eq *EventQueue) Dequeue() Pair[string, Event] {
	if len(eq.events) == 0 {
		return Pair[string, Event]{}
	}
	event := eq.events[0]
	eq.events = eq.events[1:]
	return event
}

// emit delivers e to the subscribers of its name through r.
// Pairs whose event is nil are silently skipped.
func (eq *EventQueue) emit(r *Registry, e Pair[string, Event]) {
	if e.Second != nil {
		r.EmitEventByName(e.First, e.Second)
	}
}

// IsEmpty reports whether the queue contains no events.
func (eq *EventQueue) IsEmpty() bool {
	return len(eq.events) == 0
}

// Clear discards every event still waiting in the queue without
// emitting them.
func (eq *EventQueue) Clear() {
	eq.events = make([]Pair[string, Event], 0)
}

// ProcessManyEvents dequeues and emits up to n events through r, in
// FIFO order, stopping early if the queue becomes empty.
//
// A *BreakpointEvent does not count toward the n limit: it is emitted
// like any other event but the counter is not advanced for it.
//
// Note: the returned slice currently contains only the breakpoint
// events that were encountered, not the regular events that were
// processed. This is the opposite of ProcessAllEvents and Process.
func (eq *EventQueue) ProcessManyEvents(r *Registry, n int) []Pair[string, Event] {
	processed := []Pair[string, Event]{}

	for i := 0; i < n && !eq.IsEmpty(); i++ {
		event := eq.Dequeue()
		eq.emit(r, event)
		if _, ok := event.Second.(*BreakpointEvent); ok {
			processed = append(processed, event)
			i--
		}
	}
	return processed
}

// ProcessAllEvents dequeues and emits every event in the queue through r,
// in FIFO order, until the queue is empty. Events enqueued by
// subscribers while this runs are processed too.
//
// *BreakpointEvent values are emitted but do not stop the processing,
// and they are left out of the returned slice. The slice contains all
// the other processed events, in the order they were handled.
func (eq *EventQueue) ProcessAllEvents(r *Registry) []Pair[string, Event] {
	processed := []Pair[string, Event]{}
	for !eq.IsEmpty() {
		e := eq.Dequeue()
		eq.emit(r, e)
		if _, ok := e.Second.(*BreakpointEvent); !ok {
			processed = append(processed, e)
		}
	}
	return processed
}

// Process dequeues and emits events through r in FIFO order until the
// queue is empty or a *BreakpointEvent is reached. The breakpoint is
// consumed (and emitted) but is not included in the returned slice;
// events queued after it stay in the queue for the next call.
//
// This is the method to use for per-frame processing: enqueue a
// breakpoint at the end of a frame's events, and each call to Process
// handles exactly one frame. The returned slice contains the regular
// events that were handled, in order.
func (eq *EventQueue) Process(r *Registry) []Pair[string, Event] {
	processed := []Pair[string, Event]{}
	for !eq.IsEmpty() {
		e := eq.Dequeue()
		eq.emit(r, e)
		if _, ok := e.Second.(*BreakpointEvent); ok {
			break
		} else {
			processed = append(processed, e)
		}
	}
	return processed
}
