package kv

import "encoding/json"

type Codec interface {
	// Marshal encodes a Go value to a slice of bytes.
	Marshal(v any) ([]byte, error)
	// Unmarshal decodes a slice of bytes into a Go value.
	Unmarshal(data []byte, v any) error
}

// JSON encodes/decodes Go values to/from JSON.
type JSON struct{}

// Marshal encodes a Go value to JSON.
func (c JSON) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes a JSON value into a Go value.
func (c JSON) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

var DefaultCodec = JSON{}
