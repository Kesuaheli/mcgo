package types

import (
	"fmt"
	"io"
)

const (
	SEGMENT_BITS = 0b01111111
	CONTINUE_BIT = 0b10000000
)

func ReadVarInt(r io.Reader) (i uint32, err error) {
	var pos int
	for {
		var current_byte byte
		current_byte, err = ReadOne(r)
		if err != nil {
			break
		}
		i |= (uint32(current_byte&SEGMENT_BITS) << (pos * 7))
		if current_byte&CONTINUE_BIT == 0 {
			break
		}
		pos++
		if pos > 4 {
			err = fmt.Errorf("value of VarInt is too big")
			break
		}
	}
	return
}

func PopVarInt(data *[]byte) (i uint32, err error) {
	if data == nil {
		return 0, fmt.Errorf("data for VarInt is nil")
	}
	var pos int
	for {
		if len(*data) == 0 {
			err = fmt.Errorf("not enough bytes in data to read VarInt")
			break
		}
		current_byte := (*data)[0]
		*data = (*data)[1:]
		i |= (uint32(current_byte&SEGMENT_BITS) << (pos * 7))
		if current_byte&CONTINUE_BIT == 0 {
			break
		}
		pos++
		if pos > 4 {
			err = fmt.Errorf("value of VarInt is too big")
			break
		}
	}
	return
}
