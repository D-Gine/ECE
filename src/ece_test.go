package ece

import (
	"testing"
)

func TestComponentRegistration(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Error("error: could not create registry")
	}
	RegisterComponents[int](reg)
	integers := GetComponents[int](reg)
	if integers == nil {
		t.Error("error: could not get registred component type 'int'")
	}
	nothing := GetComponents[string](reg)
	if nothing != nil {
		t.Error("error: huh ?")
	}
}

func TestAddingComponent(t *testing.T) {
	reg := NewRegistry()
	RegisterComponents[int](reg)
	integers := GetComponents[int](reg)
	AddComponent[int](reg, 0, 5)
	val, err := integers.Get(0)
	if err != nil {
		t.Errorf(`error: trying to get component from entity: %v`, err)
	}
	if val != 5 {
		t.Error("error: wrong value from getting component")
	}
}
