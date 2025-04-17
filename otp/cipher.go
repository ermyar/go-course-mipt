//go:build !solution

package otp

import (
	"io"
)

type streamCipherReader struct {
	r    io.Reader
	prng io.Reader
}

func (sc streamCipherReader) Read(p []byte) (length int, err error) {
	length, err = sc.r.Read(p)
	slice := make([]byte, length)
	sc.prng.Read(slice)
	for i := range slice {
		p[i] ^= slice[i]
	}
	return length, err
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	sc := streamCipherReader{r, prng}
	intfc := io.Reader(sc)
	return intfc
}

type streamCipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (sc streamCipherWriter) Write(p []byte) (length int, err error) {
	slice := make([]byte, len(p))
	sc.prng.Read(slice)
	for i := range p {
		slice[i] ^= p[i]
	}
	length, err = sc.w.Write(slice)
	return length, err
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	sc := streamCipherWriter{w, prng}
	i := io.Writer(sc)
	return i
}
