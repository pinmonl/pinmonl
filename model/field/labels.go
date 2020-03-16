package field

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
)

// Labels treats query string liked value as a map.
type Labels map[string]string

// Scan implements sql.Scanner interface.
func (l *Labels) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		vs := value.(string)
		if vs == "" {
			*l = nil
			return nil
		}
		q, err := url.ParseQuery(vs)
		if err != nil {
			return err
		}
		for k := range q {
			(*l)[k] = q.Get(k)
		}
	default:
	}
	return nil
}

// String returns encoded query string.
func (l Labels) String() string {
	q := url.Values{}
	for k, v := range l {
		q.Set(k, v)
	}
	return q.Encode()
}

// Value implements driver.Valuer interface.
func (l Labels) Value() (driver.Value, error) {
	return l.String(), nil
}

// MarshalJSON implements json.Marshaler interface.
func (l Labels) MarshalJSON() ([]byte, error) {
	l2 := (map[string]string)(l)
	return json.Marshal(l2)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (l *Labels) UnmarshalJSON(data []byte) error {
	v := make(map[string]string)
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*l = Labels(v)
	return nil
}

// Equal reports whether l and s contain identical name-value pair.
func (l Labels) Equal(s Labels) bool {
	return l.String() == s.String()
}
