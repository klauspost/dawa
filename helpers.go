package dawa

import (
	"io"
)

// readByteSkippingSpace() reads through an io.Reader until a character that is
// not whitespace is encountered
func readByteSkippingSpace(r io.Reader) (b byte, err error) {
	buf := make([]byte, 1)
	for {
		_, err := r.Read(buf)
		if err != nil {
			return 0, err
		}
		b := buf[0]
		switch b {
		// Only handling ASCII white space for now
		case ' ', '\t', '\n', '\v', '\f', '\r':
			continue
		default:
			return b, nil
		}
	}
}
