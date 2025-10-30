package transport

import "io"

type ReadCloser struct {
	Origin io.ReadCloser
	Encode io.ReadCloser
}

func (e *ReadCloser) Read(p []byte) (n int, err error) {
	return e.Encode.Read(p)
}

func (e *ReadCloser) Close() error {
	defer e.Origin.Close()
	return e.Encode.Close()
}
