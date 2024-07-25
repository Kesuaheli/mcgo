package types

import (
	"encoding/binary"
	"fmt"
	"io"
)

func PopLong(data *[]byte) (int64, error) {
	if data == nil || len(*data) < 8 {
		return 0, fmt.Errorf("data for Long is nil or too short")
	}

	l := binary.BigEndian.Uint64(*data)
	*data = (*data)[8:]
	return int64(l), nil
}

func WriteLong(w io.Writer, l int64) error {
	return binary.Write(w, binary.BigEndian, l)
}
