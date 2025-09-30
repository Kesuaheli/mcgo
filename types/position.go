package types

import (
	"fmt"
	"io"
)

func WriteAngle(w io.Writer, deg int) (err error) {
	_, err = w.Write([]byte{byte(deg * 256 / 360)})
	return err
}

// WritePositionXYZ is a helper function to write three Doubles to a writer.
//
// Not to be confused with [WriteBlockPosition].
func WritePositionXYZ(w io.Writer, x, y, z float64) (err error) {
	if err = WriteDouble(w, x); err != nil {
		return err
	}
	if err = WriteDouble(w, y); err != nil {
		return err
	}
	return WriteDouble(w, z)
}

// WriteBlockPosition writes a block position as Long to a writer.
//
// Not to be confused with [WritePositionXYZ].
func WriteBlockPosition(w io.Writer, x, y, z int) error {
	if x < -33554432 || x > 33554431 ||
		y < -2048 || y > 2047 ||
		z < -33554432 || z > 33554431 {
		panic(fmt.Sprintf("block position out of range: (%d, %d, %d)\n", x, y, z))
	}

	return WriteLong(w, int64(int64(x)<<26+12)|(int64(z)<<12)|int64(y))
}
