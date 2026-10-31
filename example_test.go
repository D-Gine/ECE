package ece

import (
	"fmt"
)

type Damage struct {
	Target int
	Amount int
}

type Health struct {
	Points int
}

// Example shows the whole ECE flow: components live in a SparseArray
// indexed by entity ID, there are no systems, and all behaviour is a
// callback subscribed to an event. Events can be emitted right away or
// queued and processed frame by frame using a BreakpointEvent.
func Example() {
	r := NewRegistry()

	healths := NewSparseArray[*Health]()
	healths.Insert(1, &Health{Points: 100})

	r.SubscribeToEventByName("damage", func(r *Registry, e Event) any {
		d := e.(Damage)
		if h, err := healths.Get(d.Target); err == nil {
			h.Points -= d.Amount
			fmt.Printf("entity %d took %d damage, %d HP left\n", d.Target, d.Amount, h.Points)
		}
		return nil
	}, 0)

	r.EmitEventByName("damage", Damage{Target: 1, Amount: 10})

	r.EnqueueEventByName("damage", Damage{Target: 1, Amount: 20})
	r.EnqueueEventByName("damage", Damage{Target: 1, Amount: 30})
	r.EnqueueEvent(&BreakpointEvent{})
	r.EnqueueEventByName("damage", Damage{Target: 1, Amount: 5})

	fmt.Println("frame 1:", len(r.ProcessEvents()), "events")
	fmt.Println("frame 2:", len(r.ProcessEvents()), "events")

	// Output:
	// entity 1 took 10 damage, 90 HP left
	// entity 1 took 20 damage, 70 HP left
	// entity 1 took 30 damage, 40 HP left
	// frame 1: 2 events
	// entity 1 took 5 damage, 35 HP left
	// frame 2: 1 events
}

// === SPARSE ARRAY === //

func ExampleSparseArray_Insert() {
	integers := NewSparseArray[int]()
	entity_id := 10

	integers.Insert(entity_id, 67)
	val, _ := integers.Get(entity_id)
	fmt.Println(val)

	// Output: 67
}

func ExampleSparseArray_Get() {
	integers := NewSparseArray[int]()
	entity_id := 10

	integers.Insert(entity_id, 67)
	val, _ := integers.Get(entity_id)
	fmt.Println(val)

	// Output: 67
}

func ExampleSparseArray_Remove() {
	integers := NewSparseArray[int]()
	entity_id := 10

	integers.Insert(entity_id, 67)
	err := integers.Remove(entity_id)

	_, err = integers.Get(entity_id)
	fmt.Println(err.Error())

	// Output: value is empty
}
