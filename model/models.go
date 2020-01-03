package model

import (
	"reflect"
)

// Morphable defines the interface for polymorphic relation association.
type Morphable interface {
	MorphKey() string
	MorphName() string
}

// MustBeMorphables converts slice of model to Morphable.
func MustBeMorphables(list interface{}) []Morphable {
	rl := reflect.ValueOf(list)
	if rl.Kind() != reflect.Slice {
		panic("First parameter must be slice")
	}

	out := make([]Morphable, rl.Len())
	for i := 0; i < rl.Len(); i++ {
		m, ok := rl.Index(i).Interface().(Morphable)
		if !ok {
			panic("One of the models is not morphable")
		}
		out[i] = m
	}
	return out
}
