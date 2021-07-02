package common

import (
	"testing"
)

func TestQueue_Remove(t *testing.T) {

	queue := NewQueue()
	element := queue.Remove()

	e := element.(int)
	if e != n1.element.(int) {
		t.Fatalf("Esperado %v tenemos %v", n1.element, e)
	}
	element = queue.Remove()

	e = element.(int)
	if e != n2.element.(int) {
		t.Fatalf("Esperado %v tenemos %v", n2.element, e)
	}

	element = queue.Remove()

	e = element.(int)
	if e != n3.element.(int) {
		t.Fatalf("Esperado %v tenemos %v", n3.element, e)
	}
	element = queue.Remove()

	if element != nil {
		t.Fatalf("Esperado nil")
	}
}

func TestQueue_Add(t *testing.T) {
	queue := Queue{
		size:  0,
		First: nil,
		Last:  nil,
	}

	queue.Add(1)
	queue.Add(1)
	queue.Add(1)
	queue.Add(1)

	count := 0
	for queue.HasNext() {
		queue.Remove()
		count++
	}

	if count != 4 {
		t.Fatal("Error adding")
	}

}

func TestQueue_HasNext(t *testing.T) {
	queue := Queue{
		size:  0,
		First: nil,
		Last:  nil,
	}

	if queue.HasNext() {
		t.Fatal("Shouldn't have next after creation")
	}
	queue.Add(1)

	if !queue.HasNext() {
		t.Fatal("Should have next afte add")
	}

	queue.Remove()
	if queue.HasNext() {
		t.Fatal("Shouldn't have next after remove")
	}

}

// Apoyo

func NewQueue() *Queue {
	p := Queue{
		size:  3,
		First: &n1,
		Last:  &n3,
	}
	return &p
}

var (
	n3 = node{
		element: 3,
		Next:    nil,
	}
	n2 = node{
		element: 2,
		Next:    &n3,
	}

	n1 = node{
		element: 6,
		Next:    &n2,
	}
)
