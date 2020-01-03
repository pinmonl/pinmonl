package field

import (
	"bytes"
	"testing"
	"time"
)

func TestTimeScan(t *testing.T) {
	tests := []struct {
		v     interface{}
		wants Time
	}{
		{
			v:     nil,
			wants: Time{},
		},
		{
			v:     time.Date(2019, 12, 30, 1, 2, 3, 4, time.UTC),
			wants: Time(time.Date(2019, 12, 30, 1, 2, 3, 4, time.UTC)),
		},
		{
			v:     time.Time{},
			wants: Time{},
		},
	}

	for _, test := range tests {
		got := Time{}
		err := got.Scan(test.v)
		if err != nil {
			t.Errorf("scan: error(%s)", err)
		}
		if got != test.wants {
			t.Errorf("scan: does not match")
		}
	}
}

func TestTimeString(t *testing.T) {
	tests := []struct {
		v     Time
		wants string
	}{
		{
			v:     Time{},
			wants: time.Time{}.String(),
		},
		{
			v:     Time(time.Date(2019, 12, 30, 1, 2, 3, 4, time.UTC)),
			wants: time.Date(2019, 12, 30, 1, 2, 3, 4, time.UTC).String(),
		},
	}

	for _, test := range tests {
		got := test.v.String()
		if got != test.wants {
			t.Errorf("string: expects %s, got %s", test.wants, got)
		}
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		v     Time
		wants []byte
	}{
		{
			v:     Time{},
			wants: []byte("null"),
		},
		{
			v:     Time(time.Date(2019, 12, 30, 1, 2, 3, 4, time.UTC)),
			wants: []byte(`"2019-12-30T01:02:03.000000004Z"`),
		},
	}

	for _, test := range tests {
		got, err := test.v.MarshalJSON()
		if err != nil {
			t.Errorf("json: marshal error(%s)", err)
		}
		if bytes.Compare(got, test.wants) != 0 {
			t.Errorf("json: marshal expects %s, got %s", test.wants, got)
		}
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	tests := []struct{
		v []byte
		wants Time
	}{
		{
			v: []byte("null"),
			wants: Time{},
		},
		{
			v: []byte(`"2019-12-30T01:02:03.000000004Z"`),
			wants: Time(time.Date(2019, 12, 30, 1, 2, 3, 4, time.UTC)),
		},
	}

	for _, test := range tests {
		got := Time{}
		err := got.UnmarshalJSON(test.v)
		if err != nil {
			t.Errorf("json: unmarshal error(%s)", err)
		}
		if got != test.wants {
			t.Errorf("json: unmarshal does not match of %q", test.v)
		}
	}
}
