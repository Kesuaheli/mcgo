package types

import (
	"fmt"
	"io"
)

const (
	SEGMENT_BITS = 0b01111111
	CONTINUE_BIT = 0b10000000
)

func ReadVarInt(r io.Reader) (i int32, err error) {
	var pos int
	for {
		var current_byte byte
		current_byte, err = ReadOne(r)
		if err != nil {
			break
		}
		i |= (int32(current_byte&SEGMENT_BITS) << (pos * 7))
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

func WriteVarInt(w io.Writer, i int32) (err error) {
	for {
		if i & ^int32(SEGMENT_BITS) == 0 {
			_, err = w.Write([]byte{byte(i)})
			return
		}
		b := byte(i&int32(SEGMENT_BITS) | CONTINUE_BIT)
		w.Write([]byte{b})

		i >>= 7
	}
}

func PopVarInt(data *[]byte) (i int32, err error) {
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
		i |= (int32(current_byte&SEGMENT_BITS) << (pos * 7))
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
