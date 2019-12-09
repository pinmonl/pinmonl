package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

// ParseJSON reads bytes into interface
func ParseJSON(in []byte, out interface{}) error {
	err := json.Unmarshal(in, &out)
	return err
}

// ToJSON stringifies input to bytes
func ToJSON(in interface{}) ([]byte, error) {
	out, err := json.Marshal(in)
	return out, err
}

// ReaderToByte converts io.Reader to bytes
func ReaderToByte(reader io.Reader) []byte {
	b := &bytes.Buffer{}
	_, err := b.ReadFrom(reader)
	if err != nil {
		return nil
	}
	return b.Bytes()
}
