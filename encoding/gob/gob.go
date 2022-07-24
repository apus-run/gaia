package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/apus-run/gaia/encoding"
)

// Name is the name registered for the gob codec.
const Name = "gob"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with gob.
type codec struct{}

// Marshal gob encode
func (codec) Marshal(v interface{}) ([]byte, error) {
	var (
		buffer bytes.Buffer
	)

	err := gob.NewEncoder(&buffer).Encode(v)
	return buffer.Bytes(), err
}

// Unmarshal gob encode
func (codec) Unmarshal(data []byte, value interface{}) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(value)
	if err != nil {
		return err
	}
	return nil
}

func (codec) Name() string {
	return Name
}
