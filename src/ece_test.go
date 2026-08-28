package ece

import "fmt"

func main() {
	reg := NewRegistry()
	RegisterComponents[int](reg)
	integers := GetComponents[int](reg)
	AddComponent[int](reg, 0, 5)
	fmt.Println(integers.Get(0))
}
