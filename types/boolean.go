package types

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadBoolean(r io.Reader) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("data for boolean is nil or too short")
	}
	l, err := ReadOne(r)
	return l != 0, err
}

func WriteBoolean(w io.Writer, l bool) error {
	return binary.Write(w, binary.BigEndian, l)
}
