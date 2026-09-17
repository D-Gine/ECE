package ece

import (
	"reflect"
	"testing"
)

const test_entity int = 10

func TestComponentRegistration(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Error("error: could not create registry")
	}
	reg.RegisterComponents[int]()
	integers := reg.GetComponents[int]()
	if integers == nil {
		t.Error("error: could not get registred component type 'int'")
	}
	nothing := reg.GetComponents[string]()
	if nothing != nil {
		t.Error("error: huh ?")
	}
}

func TestAddingComponent(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterComponents[int]()
	integers := reg.GetComponents[int]()

	reg.AddComponent(test_entity, 5)
	val, err := integers.Get(test_entity)
	if err != nil {
		t.Errorf(`error: trying to get component from entity: %v`, err)
	}
	if val != 5 {
		t.Error("error: wrong value from getting component")
	}
}

func TestRemovingComponent(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterComponents[int]()
	integers := reg.GetComponents[int]()

	reg.AddComponent(test_entity, 5)
	_, err := integers.Get(test_entity)
	if err != nil {
		t.Errorf(`error: Get component not working: can't test remover`)
	}
	err = integers.Remove(test_entity)
	if err != nil {
		t.Errorf(`error: could not remove component: %v`, err)
	}
	_, err = integers.Get(test_entity)
	if err == nil {
		t.Error("error: component not removed")
	}
}

func TestEventSubscription(t *testing.T) {
	reg := NewRegistry()

	reg.SubscribeToEvent(func(r *Registry, e *CoucouEvent) any {
		return nil
	}, 1)

	if reg.events_subscribers[reflect.TypeOf(&CoucouEvent{})] == nil {
		t.Error("error: could not subscribe to event")
	}
}

func TestEventTriggering(t *testing.T) {
	reg := NewRegistry()
	var res string

	reg.SubscribeToEvent(func(r *Registry, e *CoucouEvent) any {
		res = e.text
		return nil
	}, 1)

	e := &CoucouEvent{text: "Hello, World!"}
	reg.eventQueue.Enqueue(e)
	reg.eventQueue.Process(reg)
	if res != "Hello, World!" {
		t.Errorf("error: event not triggered correctly, got '%s'", res)
	}
}

func TestEventQueueProcessingWithBreakpoint(t *testing.T) {
	reg := NewRegistry()
	var res string

	reg.SubscribeToEvent(func(r *Registry, e *CoucouEvent) any {
		res += e.text
		return nil
	}, 1)

	reg.SubscribeToEvent(func(r *Registry, e *BreakpointEvent) any {
		res += "Breakpoint!"
		return nil
	}, 1)

	e1 := &CoucouEvent{text: "Hello, "}
	e2 := &BreakpointEvent{}
	e3 := &CoucouEvent{text: "World!"}

	reg.eventQueue.Enqueue(e1)
	reg.eventQueue.Enqueue(e2)
	reg.eventQueue.Enqueue(e3)

	reg.eventQueue.Process(reg)

	if res != "Hello, Breakpoint!" {
		t.Errorf("error: event queue processing with breakpoint failed, got '%s'", res)
	}
}
