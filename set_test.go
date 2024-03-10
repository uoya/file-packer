package main

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	set := NewSet()
	set.Add("item1")
	set.Add("item1")
	set.Add("item1")
	set.Add("item2")
	set.Add("item2")
	set.Add("item3")

	expect := make(Set)
	expect["item1"] = true
	expect["item2"] = true
	expect["item3"] = true
	if !reflect.DeepEqual(expect, set) {
		t.Errorf("expect %v, got %v", expect, set)
	}
}

func TestValues(t *testing.T) {
	set := NewSet()
	set.Add("item1")
	set.Add("item2")
	set.Add("item3")

	expect := []string{"item1", "item2", "item3"}
	if !reflect.DeepEqual(expect, set.Values()) {
		t.Errorf("expect %v, got %v", expect, set.Values())
	}
}
