package helpers

import "io"

// Fake Tool for enabling io.ReadCloser
type NopCloser struct {
	io.Reader
}

func (NopCloser) Close() error { return nil }
