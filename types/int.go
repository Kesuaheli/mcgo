package types

import (
	"encoding/binary"
	"fmt"
	"io"
)

func PopUInt(data *[]byte) (uint32, error) {
	if data == nil || len(*data) < 4 {
		return 0, fmt.Errorf("data for Unsigned Int is nil or too short")
	}

	i := binary.BigEndian.Uint32(*data)
	*data = (*data)[4:]
	return i, nil
}

func PopInt(data *[]byte) (int32, error) {
	i, err := PopUInt(data)
	return int32(i), err
}

func WriteInt[I uint32 | int32](w io.Writer, i I) error {
	return binary.Write(w, binary.BigEndian, i)
}
