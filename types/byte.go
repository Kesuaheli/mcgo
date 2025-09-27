package types

import (
	"fmt"
)

func PopUByte(data *[]byte) (uint8, error) {
	if data == nil || len(*data) < 1 {
		return 0, fmt.Errorf("data for Unsigned Byte is nil or too short")
	}

	i := (*data)[0]
	*data = (*data)[1:]
	return i, nil
}

func PopByte(data *[]byte) (int8, error) {
	val, err := PopUByte(data)
	return int8(val), err
}
