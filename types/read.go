package types

import "io"

func ReadOne(r io.Reader) (b byte, err error) {
	var p [1]byte
	_, err = r.Read(p[:])
	return p[0], err
}

func Read(n int, r io.Reader) (b []byte, err error) {
	b = make([]byte, n)
	n, err = r.Read(b)
	return b[:n], err
}
