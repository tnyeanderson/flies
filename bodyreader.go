package flies

import (
	"bytes"
	"io"
)

// BodyReader Allows "rereading" of the content of the body.  On first read, it
// saves the entire body in memory. Then every time [BodyReader.Close] is
// called, a new internal reader is created based on the cached body.
//
// A BodyReader can be used as a drop-in replacement for an
// [http.Request.Body].
type BodyReader struct {
	content []byte
	buffer  *bytes.Buffer
	reader  io.Reader
}

// NewBodyReader returns an initialized [BodyReader] based on the provided
// [io.Reader].
func NewBodyReader(r io.Reader) *BodyReader {
	b := &BodyReader{}
	b.buffer = &bytes.Buffer{}
	b.reader = io.TeeReader(r, b.buffer)
	return b
}

// Read reads the body content from the internal reader.
func (b *BodyReader) Read(p []byte) (int, error) {
	return b.reader.Read(p)
}

// Close saves the content of the buffered body, and resets the internal reader
// to "start over" from the beginning of the content on the next call to
// [BodyReader.Read].
func (b *BodyReader) Close() error {
	if b.content == nil {
		b.content = b.buffer.Bytes()
	}
	b.reader = bytes.NewReader(b.content)
	return nil
}
