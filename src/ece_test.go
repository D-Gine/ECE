package ece

import (
	"testing"
)

const test_entity int = 10

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

	AddComponent[int](reg, test_entity, 5)
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
	RegisterComponents[int](reg)
	integers := GetComponents[int](reg)

	AddComponent(reg, test_entity, 5)
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

type HitEntity struct {
	e int
}

func onCoucouEventEncore(r *Registry, e *CoucouEvent) {
	println("Weeee alors on a encore reçu un event coucou", e.text)
}
func TestEventCreation(t *testing.T) {
	reg := NewRegistry()

	subscribe(reg, onCoucouEvent, 0)
	subscribe(reg, onCoucouEventEncore, 1)

	e := &CoucouEvent{text: "Hello, World!"}
	e.triggerEvent(reg)
}
