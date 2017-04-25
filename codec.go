package queue

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// Codec define interface for codecs to encode/decode objects in queues
type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// JSONCodec provides JSON based Codec
type JSONCodec byte

// Marshal returns the JSON encoding of v.
func (c *JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func (c *JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewJSONCodec creates new JSONCodec
func NewJSONCodec() *JSONCodec {
	r := JSONCodec(0)
	return &r
}

// GOBCodec provides GOB based Codec
type GOBCodec byte

// Marshal returns the GOB encoding of v.
func (c *GOBCodec) Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)
	return b.Bytes(), err
}

// Unmarshal parses the GOB-encoded data and stores the result
// in the value pointed to by v.
func (c *GOBCodec) Unmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}

// NewGOBCodec creates new GOBCodec
func NewGOBCodec() *GOBCodec {
	r := GOBCodec(0)
	return &r
}
