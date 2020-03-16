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

// MorphableList is slice of Morphable.
type MorphableList []Morphable

// Keys extracts key from each Morphable.
func (ml MorphableList) Keys() []string {
	out := make([]string, len(ml))
	for i, m := range ml {
		out[i] = m.MorphKey()
	}
	return out
}

// Names extracts name from each Morphable.
func (ml MorphableList) Names() []string {
	checks := map[string]bool{}
	out := []string{}
	for _, m := range ml {
		k := m.MorphName()
		if has, ok := checks[k]; !ok || !has {
			checks[k] = true
			out = append(out, k)
		}
	}
	return out
}
