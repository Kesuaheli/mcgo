package types

import (
	"encoding/binary"
	"fmt"
	"io"
)

func PopDouble(data *[]byte) (float64, error) {
	if data == nil || len(*data) < 8 {
		return 0, fmt.Errorf("data for Long is nil or too short")
	}

	l := binary.BigEndian.Uint64(*data)
	*data = (*data)[8:]
	return float64(l), nil
}

func WriteDouble(w io.Writer, l float64) error {
	return binary.Write(w, binary.BigEndian, l)
}

func PopFloat(data *[]byte) (float32, error) {
	if data == nil || len(*data) < 4 {
		return 0, fmt.Errorf("data for Long is nil or too short")
	}

	l := binary.BigEndian.Uint32(*data)
	*data = (*data)[4:]
	return float32(l), nil
}

func WriteFloat(w io.Writer, l float32) error {
	return binary.Write(w, binary.BigEndian, l)
}
