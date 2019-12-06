package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func ParseJsonFrom(in []byte, out interface{}) error {
	err := json.Unmarshal(in, &out)
	return err
}

func RenderAsJson(in interface{}) ([]byte, error) {
	out, err := json.Marshal(in)
	return out, err
}

func ReaderToByte(reader io.Reader) []byte {
	b := &bytes.Buffer{}
	_, err := b.ReadFrom(reader)
	if err != nil {
		return nil
	}
	return b.Bytes()
}
