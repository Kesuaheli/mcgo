package types

import (
	"encoding/binary"
	"fmt"
)

func PopUShort(data *[]byte) (uint16, error) {
	if data == nil || len(*data) < 2 {
		return 0, fmt.Errorf("data for Unsigned Short is nil or too short")
	}

	i := binary.BigEndian.Uint16(*data)
	*data = (*data)[2:]
	return i, nil
}

func PopShort(data *[]byte) (int16, error) {
	val, err := PopUShort(data)
	return int16(val), err
}
