package field

import (
	"bytes"
	"testing"
)

func TestLabelsString(t *testing.T) {
	tests := []struct {
		v     Labels
		wants string
	}{
		{
			v:     Labels{"a": "1", "b": "2", "c": "3", "d": "4"},
			wants: "a=1&b=2&c=3&d=4",
		},
		{
			v:     Labels{"\"awe\"#![]": "escaped"},
			wants: "%22awe%22%23%21%5B%5D=escaped",
		},
	}

	for _, test := range tests {
		got := test.v.String()
		if got != test.wants {
			t.Errorf("Labels expects %q, got %q", test.wants, got)
		}
	}
}

func TestLabelScan(t *testing.T) {
	tests := []struct {
		v     interface{}
		wants Labels
	}{
		{
			v:     "a=1&b=2&c=3",
			wants: Labels{"a": "1", "b": "2", "c": "3"},
		},
		{
			v:     "q=select+%2A+from+table",
			wants: Labels{"q": "select * from table"},
		},
	}

	for _, test := range tests {
		lb := Labels{}
		err := lb.Scan(test.v)
		if err != nil {
			t.Errorf("scan: error(%s)", err)
		}
		for k, wants := range test.wants {
			got := lb[k]
			if got != wants {
				t.Errorf("scan: expects %q, got %q", wants, got)
			}
		}
	}
}

func TestLabelsMarshalJSON(t *testing.T) {
	tests := []struct {
		v     Labels
		wants []byte
	}{
		{
			v:     Labels{"a": "1", "b": "2"},
			wants: []byte(`{"a":"1","b":"2"}`),
		},
	}

	for _, test := range tests {
		got, err := test.v.MarshalJSON()
		if err != nil {
			t.Errorf("json: marshal error(%s)", err)
		}
		if bytes.Compare(got, test.wants) != 0 {
			t.Errorf("json: marshal expects %q, got %q", test.wants, got)
		}
	}
}

func TestLabelsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		v     []byte
		wants Labels
	}{
		{
			v:     []byte(`{"a": "1", "b": "2"}`),
			wants: Labels{"a": "1", "b": "2"},
		},
		{
			v:     []byte(`{"q": "select * from table", "b": "b"}`),
			wants: Labels{"q": "select * from table", "b": "b"},
		},
	}

	for _, test := range tests {
		lb := Labels{}
		err := lb.UnmarshalJSON(test.v)
		if err != nil {
			t.Errorf("json: unmarshal error(%s)", err)
		}
		for k, wants := range test.wants {
			got := lb[k]
			if got != wants {
				t.Errorf("json: unmarshal expects %s, got %s", wants, got)
			}
		}
	}
}
