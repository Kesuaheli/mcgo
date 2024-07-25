package types

import (
	"fmt"
	"io"
	"math"
)

func WriteString(w io.Writer, s string) (err error) {
	return WriteStringData(w, []byte(s))
}

func WriteStringData(w io.Writer, data []byte) (err error) {
	if len(data) > math.MaxInt32 {
		return fmt.Errorf("write string: string is too long: %d/%d", len(data), math.MaxInt32)
	}
	err = WriteVarInt(w, int32(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func PopString(data *[]byte) (s string, err error) {
	length, err := PopVarInt(data)
	if err != nil {
		return
	}
	s = string((*data)[:length])
	*data = (*data)[length:]
	return
}
