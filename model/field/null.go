package field

import "database/sql"

type NullBool sql.NullBool

func NewNullBool(v bool) NullBool {
	b := NullBool{}
	b.SetValid(v)
	return b
}

func (b *NullBool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

func (b NullBool) Value() bool {
	return b.Bool
}

type NullInt64 sql.NullInt64

func NewNullInt64(v int64) NullInt64 {
	i := NullInt64{}
	i.SetValid(v)
	return i
}

func (i *NullInt64) SetValid(v int64) {
	i.Int64 = v
	i.Valid = true
}

func (i NullInt64) Value() int64 {
	return i.Int64
}

type NullFloat64 sql.NullFloat64

func NewNullFloat64(v float64) NullFloat64 {
	f := NullFloat64{}
	f.SetValid(v)
	return f
}

func (f *NullFloat64) SetValid(v float64) {
	f.Float64 = v
	f.Valid = true
}

func (f NullFloat64) Value() float64 {
	return f.Float64
}

type NullString sql.NullString

func NewNullString(v string) NullString {
	s := NullString{}
	s.SetValid(v)
	return s
}

func (s *NullString) SetValid(v string) {
	s.String = v
	s.Valid = true
}

func (s NullString) Value() string {
	return s.String
}

type NullValue struct {
	Ref   interface{}
	Valid bool
}

func NewNullValue(v interface{}) NullValue {
	n := NullValue{}
	n.SetValid(v)
	return n
}

func (n *NullValue) SetValid(v interface{}) {
	n.Ref = v
	n.Valid = true
}

func (n NullValue) Value() interface{} {
	return n.Ref
}
