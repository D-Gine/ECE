package ece

import "fmt"

func Example() {
	fmt.Println("1312")
	// Output: 1312
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
	integers.Remove(entity_id)

	_, err := integers.Get(entity_id)
	fmt.Println(err.Error())

	// Output: value is empty
}
