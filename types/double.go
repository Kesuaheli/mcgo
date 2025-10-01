package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func PopDouble(data *[]byte) (parsedDouble float64, err error) {
	if data == nil || len(*data) < 8 {
		return 0, fmt.Errorf("data for Double is nil or too short")
	}

	r := bytes.NewBuffer(*data)
	err = binary.Read(r, binary.BigEndian, &parsedDouble)
	*data = (*data)[8:]
	return parsedDouble, err
}

func WriteDouble(w io.Writer, l float64) error {
	return binary.Write(w, binary.BigEndian, l)
}

func PopFloat(data *[]byte) (parsedFloat float32, err error) {
	if data == nil || len(*data) < 4 {
		return 0, fmt.Errorf("data for Float is nil or too short")
	}

	r := bytes.NewBuffer(*data)
	err = binary.Read(r, binary.BigEndian, &parsedFloat)
	*data = (*data)[4:]
	return parsedFloat, err
}

func WriteFloat(w io.Writer, l float32) error {
	return binary.Write(w, binary.BigEndian, l)
}
