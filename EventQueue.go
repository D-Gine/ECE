package ece

type Pair[T string, U any] struct {
	First  T
	Second U
}

type EventQueue struct {
	events []Pair[string, Event]
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		events: make([]Pair[string, Event], 0),
	}
}

func (eq *EventQueue) Enqueue(event Event) {
	eq.EnqueueByName(eventName(event), event)
}

func (eq *EventQueue) EnqueueByName(name string, event Event) {
	eq.events = append(eq.events, Pair[string, Event]{name, event})
}

func (eq *EventQueue) Dequeue() Pair[string, Event] {
	if len(eq.events) == 0 {
		return Pair[string, Event]{}
	}
	event := eq.events[0]
	eq.events = eq.events[1:]
	return event
}

func (eq *EventQueue) emit(r *Registry, e Pair[string, Event]) {
	if e.Second != nil {
		r.EmitEventByName(e.First, e.Second)
	}
}

func (eq *EventQueue) IsEmpty() bool {
	return len(eq.events) == 0
}

func (eq *EventQueue) Clear() {
	eq.events = make([]Pair[string, Event], 0)
}

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

func (eq *EventQueue) GetQueue() []Pair[string, Event] {
	return eq.events
}
