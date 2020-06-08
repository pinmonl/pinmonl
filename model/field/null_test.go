package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullBool(t *testing.T) {
	tests := []struct {
		nullable NullBool
		value    bool
		valid    bool
	}{
		{
			nullable: NullBool{},
			value:    false,
			valid:    false,
		},
		{
			nullable: NewNullBool(false),
			value:    false,
			valid:    true,
		},
		{
			nullable: NewNullBool(true),
			value:    true,
			valid:    true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.value, test.nullable.Value())
		assert.Equal(t, test.valid, test.nullable.Valid)
	}
}

func TestNullInt64(t *testing.T) {
	tests := []struct {
		nullable NullInt64
		value    int64
		valid    bool
	}{
		{
			nullable: NullInt64{},
			value:    0,
			valid:    false,
		},
		{
			nullable: NewNullInt64(0),
			value:    0,
			valid:    true,
		},
		{
			nullable: NewNullInt64(64),
			value:    64,
			valid:    true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.value, test.nullable.Value())
		assert.Equal(t, test.valid, test.nullable.Valid)
	}
}

func TestNullFloat64(t *testing.T) {
	tests := []struct {
		nullable NullFloat64
		value    float64
		valid    bool
	}{
		{
			nullable: NullFloat64{},
			value:    0,
			valid:    false,
		},
		{
			nullable: NewNullFloat64(0),
			value:    0,
			valid:    true,
		},
		{
			nullable: NewNullFloat64(64),
			value:    64,
			valid:    true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.value, test.nullable.Value())
		assert.Equal(t, test.valid, test.nullable.Valid)
	}
}

func TestNullString(t *testing.T) {
	tests := []struct {
		nullable NullString
		value    string
		valid    bool
	}{
		{
			nullable: NullString{},
			value:    "",
			valid:    false,
		},
		{
			nullable: NewNullString(""),
			value:    "",
			valid:    true,
		},
		{
			nullable: NewNullString("string"),
			value:    "string",
			valid:    true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.value, test.nullable.Value())
		assert.Equal(t, test.valid, test.nullable.Valid)
	}
}

func TestNullValue(t *testing.T) {
	tests := []struct {
		nullable NullValue
		value    interface{}
		valid    bool
	}{
		{
			nullable: NullValue{},
			value:    nil,
			valid:    false,
		},
		{
			nullable: NewNullValue(nil),
			value:    nil,
			valid:    true,
		},
		{
			nullable: NewNullValue(struct{}{}),
			value:    struct{}{},
			valid:    true,
		},
		{
			nullable: NewNullValue(64),
			value:    64,
			valid:    true,
		},
		{
			nullable: NewNullValue("string"),
			value:    "string",
			valid:    true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.value, test.nullable.Value())
		assert.Equal(t, test.valid, test.nullable.Valid)
	}
}
