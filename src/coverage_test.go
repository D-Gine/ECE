package ece

import (
	"reflect"
	"testing"
)

type coverageComponent struct {
	name string
}

type secondCoverageComponent struct {
	value int
}

type coverageEvent struct {
	value string
}

func TestOptionalSomeAndNone(t *testing.T) {
	some := Some(42)
	value, err := some.Get()
	if err != nil {
		t.Fatalf("Some.Get() returned an error: %v", err)
	}
	if value != 42 {
		t.Fatalf("Some.Get() = %d, want 42", value)
	}

	none := None[int]()
	value, err = none.Get()
	if err == nil {
		t.Fatal("None.Get() returned nil error")
	}
	if value != 0 {
		t.Fatalf("None.Get() = %d, want zero value", value)
	}
}

func TestSparseArrayGrowthAndLookup(t *testing.T) {
	s := NewSparseArray[string]()
	for entity, value := range map[int]string{
		0:  "zero",
		3:  "three",
		4:  "four",
		10: "ten",
	} {
		s.Insert(entity, value)
	}

	for entity, want := range map[int]string{
		0:  "zero",
		3:  "three",
		4:  "four",
		10: "ten",
	} {
		got, err := s.Get(entity)
		if err != nil {
			t.Errorf("Get(%d) returned an error: %v", entity, err)
		} else if got != want {
			t.Errorf("Get(%d) = %q, want %q", entity, got, want)
		}
	}

	if _, err := s.Get(1); err == nil {
		t.Error("Get(1) returned nil error for an absent component")
	}
}

func TestSparseArrayRemoveKeepsRemainingComponents(t *testing.T) {
	s := NewSparseArray[string]()
	s.Insert(1, "one")
	s.Insert(2, "two")
	s.Insert(3, "three")

	if err := s.Remove(2); err != nil {
		t.Fatalf("Remove(2) returned an error: %v", err)
	}
	if _, err := s.Get(2); err == nil {
		t.Error("Get(2) returned a value after removal")
	}
	for entity, want := range map[int]string{1: "one", 3: "three"} {
		got, err := s.Get(entity)
		if err != nil {
			t.Errorf("Get(%d) returned an error after removing another entity: %v", entity, err)
		} else if got != want {
			t.Errorf("Get(%d) = %q, want %q", entity, got, want)
		}
	}

	if err := s.Remove(2); err != nil {
		t.Fatalf("removing an absent component returned an error: %v", err)
	}
}

func TestSparseArrayRemoveLastComponent(t *testing.T) {
	s := NewSparseArray[int]()
	s.Insert(8, 99)
	if err := s.Remove(8); err != nil {
		t.Fatalf("Remove(8) returned an error: %v", err)
	}
	if len(s.dense) != 0 {
		t.Fatalf("dense length = %d, want 0", len(s.dense))
	}
	if _, err := s.Get(8); err == nil {
		t.Error("Get(8) returned a value after removing the only component")
	}
}

func TestSparseArrayRemoveRejectsCorruptDenseIndex(t *testing.T) {
	s := NewSparseArray[int]()
	s.Insert(1, 99)
	s.sparse[1] = Some(10)
	if err := s.Remove(1); err == nil {
		t.Fatal("Remove() returned nil for an invalid dense index")
	}
}

func TestSparseArrayRemoveAbsentEntityOutsideSparseRange(t *testing.T) {
	s := NewSparseArray[int]()
	s.Insert(1, 99)
	for _, entity := range []int{-1, 20} {
		if err := s.Remove(entity); err != nil {
			t.Errorf("Remove(%d) returned an error for an absent entity: %v", entity, err)
		}
	}
}

func TestEventQueueOperations(t *testing.T) {
	queue := NewEventQueue()
	if !queue.IsEmpty() {
		t.Fatal("new event queue is not empty")
	}
	if got := queue.Dequeue(); got != nil {
		t.Fatalf("Dequeue() on an empty queue = %#v, want nil", got)
	}

	first := &coverageEvent{value: "first"}
	second := &coverageEvent{value: "second"}
	queue.Enqueue(first)
	queue.Enqueue(second)
	if queue.IsEmpty() {
		t.Fatal("queue is empty after enqueue")
	}
	if got := queue.Dequeue(); got != first {
		t.Fatalf("first Dequeue() = %#v, want first event", got)
	}
	if got := queue.Dequeue(); got != second {
		t.Fatalf("second Dequeue() = %#v, want second event", got)
	}
	queue.Enqueue(first)
	queue.Clear()
	if !queue.IsEmpty() {
		t.Fatal("queue is not empty after Clear()")
	}
}

func TestEventSubscriptionPriorityAndStableOrdering(t *testing.T) {
	reg := NewRegistry()
	order := make([]string, 0, 3)
	subscribe := func(name string) {
		reg.SubscribeToEvent(func(*Registry, *coverageEvent) any {
			order = append(order, name)
			return nil
		}, 2)
	}
	subscribe("first")
	subscribe("second")
	reg.SubscribeToEvent(func(*Registry, *coverageEvent) any {
		order = append(order, "high")
		return nil
	}, 1)

	reg.EmitEvent(&coverageEvent{value: "event"})
	want := []string{"high", "first", "second"}
	if !reflect.DeepEqual(order, want) {
		t.Fatalf("subscriber order = %v, want %v", order, want)
	}
}

func TestEmitEventIgnoresUnsubscribedTypes(t *testing.T) {
	reg := NewRegistry()
	called := false
	reg.SubscribeToEvent(func(*Registry, *coverageEvent) any {
		called = true
		return nil
	}, 0)

	reg.EmitEvent(&BreakpointEvent{})
	if called {
		t.Error("subscriber for another event type was called")
	}
}

func TestProcessManyEventsHonorsBreakpointWithoutCountingIt(t *testing.T) {
	reg := NewRegistry()
	processed := make([]string, 0, 3)
	reg.SubscribeToEvent(func(_ *Registry, event *coverageEvent) any {
		processed = append(processed, "event:"+event.value)
		return nil
	}, 0)
	reg.SubscribeToEvent(func(*Registry, *BreakpointEvent) any {
		processed = append(processed, "breakpoint")
		return nil
	}, 0)

	reg.EnqueueEvent(&BreakpointEvent{})
	reg.EnqueueEvent(&coverageEvent{value: "one"})
	reg.EnqueueEvent(&coverageEvent{value: "two"})
	reg.EnqueueEvent(&coverageEvent{value: "three"})
	reg.ProcessManyEvents(2)

	if want := []string{"breakpoint", "event:one", "event:two"}; !reflect.DeepEqual(processed, want) {
		t.Fatalf("processed events = %v, want %v", processed, want)
	}
}

func TestProcessAllEventsAndProcessBreakpoint(t *testing.T) {
	reg := NewRegistry()
	processed := make([]string, 0, 3)
	reg.SubscribeToEvent(func(*Registry, *coverageEvent) any {
		processed = append(processed, "event")
		return nil
	}, 0)
	reg.SubscribeToEvent(func(*Registry, *BreakpointEvent) any {
		processed = append(processed, "breakpoint")
		return nil
	}, 0)

	reg.EnqueueEvent(&coverageEvent{})
	reg.EnqueueEvent(&BreakpointEvent{})
	reg.EnqueueEvent(&coverageEvent{})
	reg.ProcessEvents()
	if want := []string{"event", "breakpoint"}; !reflect.DeepEqual(processed, want) {
		t.Fatalf("Process() = %v, want %v", processed, want)
	}

	reg.ProcessAllEvents()
	if want := []string{"event", "breakpoint", "event"}; !reflect.DeepEqual(processed, want) {
		t.Fatalf("ProcessAllEvents() = %v, want %v", processed, want)
	}
}

func TestRegistryComponentRegistrationIsIdempotent(t *testing.T) {
	reg := NewRegistry()
	first := reg.RegisterComponents[coverageComponent]()
	first.Insert(1, coverageComponent{name: "kept"})
	second := reg.RegisterComponents[coverageComponent]()
	if first != second {
		t.Fatal("registering the same component type returned a different array")
	}
	value, err := second.Get(1)
	if err != nil || value.name != "kept" {
		t.Fatalf("component data was not preserved: value=%+v, err=%v", value, err)
	}
}

func TestRegistryComponentsAndErrors(t *testing.T) {
	reg := NewRegistry()
	if got := reg.GetComponents[coverageComponent](); got != nil {
		t.Fatal("GetComponents() returned an array for an unregistered type")
	}
	if err := reg.AddComponent(1, coverageComponent{name: "missing"}); err == nil {
		t.Fatal("AddComponent() returned nil for an unregistered type")
	}

	reg.RegisterComponents[coverageComponent]()
	reg.RegisterComponents[secondCoverageComponent]()
	if err := reg.AddComponent(4, coverageComponent{name: "ok"}); err != nil {
		t.Fatalf("AddComponent() returned an error: %v", err)
	}
	if err := reg.AddComponent(4, secondCoverageComponent{value: 7}); err != nil {
		t.Fatalf("AddComponent() returned an error for a second type: %v", err)
	}

	reg.components[reflect.TypeOf((*coverageComponent)(nil)).Elem()] = NewSparseArray[int]()
	if err := reg.AddComponent(4, coverageComponent{name: "wrong storage"}); err == nil {
		t.Fatal("AddComponent() returned nil for an incompatible stored array")
	}
}
