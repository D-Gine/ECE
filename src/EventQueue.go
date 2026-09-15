package ece

type EventQueue struct {
	events []Event
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		events: make([]Event, 0),
	}
}

func (eq *EventQueue) Enqueue(event Event) {
	eq.events = append(eq.events, event)
}

func (eq *EventQueue) Dequeue() Event {
	if len(eq.events) == 0 {
		return nil
	}
	event := eq.events[0]
	eq.events = eq.events[1:]
	return event
}

func (eq *EventQueue) IsEmpty() bool {
	return len(eq.events) == 0
}

func (eq *EventQueue) Clear() {
	eq.events = make([]Event, 0)
}

func (eq *EventQueue) ProcessManyEvents(r *Registry, n int) {
	for i := 0; i < n && !eq.IsEmpty(); i++ {
		event := eq.Dequeue()
		if event != nil {
			event.triggerEvent(r)
		}
		if _, ok := event.(*BreakpointEvent); ok {
			i--
		}
	}
}

func (eq *EventQueue) ProcessAllEvents(r *Registry) {
	for !eq.IsEmpty() {
		event := eq.Dequeue()
		if event != nil {
			event.triggerEvent(r)
		}
	}
}

func (eq *EventQueue) Process(r *Registry) {
	for !eq.IsEmpty() {
		event := eq.Dequeue()
		if event != nil {
			event.triggerEvent(r)
			if _, ok := event.(*BreakpointEvent); ok {
				break
			}
		}
	}
}
