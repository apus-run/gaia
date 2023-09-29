package ws

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/apus-run/gaia/encoding"
)

func Marshal(codec encoding.Codec, msg Any) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}

	if codec != nil {
		dataBuffer, err := codec.Marshal(msg)
		if err != nil {
			return nil, err
		}
		return dataBuffer, nil
	} else {
		switch t := msg.(type) {
		case []byte:
			return t, nil
		case string:
			return []byte(t), nil
		default:
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(msg); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
	}
}

func Unmarshal(codec encoding.Codec, buf []byte, data interface{}) error {
	if codec != nil {
		if err := codec.Unmarshal(buf, data); err != nil {
			return err
		}
	} else {
		data = buf
	}
	return nil
}
