package field

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Time wraps time.Time for database nullable timestamp.
type Time time.Time

func Now() Time {
	t := time.Now().Round(time.Second).UTC()
	return Time(t)
}

// Scan implements sql.Scanner interface.
func (t *Time) Scan(src interface{}) error {
	st, ok := src.(time.Time)
	if ok {
		*t = Time(st.UTC())
	}
	return nil
}

// Strings implements fmt.Stringer interface.
func (t Time) String() string {
	return t.Time().String()
}

// Time returns time.Time.
func (t Time) Time() time.Time {
	return (time.Time)(t)
}

// Value implements driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	t2 := t.Time()
	if t2.IsZero() {
		return nil, nil
	}
	return t2, nil
}

// MarshalJSON implements json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	t2 := t.Time()
	if t2.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t2)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	var t2 *time.Time
	err := json.Unmarshal(data, &t2)
	if err != nil {
		return err
	}
	if t2 != nil {
		*t = Time((*t2).UTC())
	}
	return nil
}
