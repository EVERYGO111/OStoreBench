package utils

import (
	"io"
	"bufio"
)

// a fake reader to generate content to upload
type FakeReader struct {
	id          uint64 // an id number
	size        int64  // max bytes
	sharedBytes []byte
}

func NewFakeReader(id uint64, fileSize int64) *FakeReader {
	return &FakeReader{id: id, size: fileSize, sharedBytes: make([]byte, 1024)}
}

func (l *FakeReader) Read(p []byte) (n int, err error) {
	if l.size <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.size {
		n = int(l.size)
	} else {
		n = len(p)
	}
	if n >= 8 {
		for i := 0; i < 8; i++ {
			p[i] = byte(l.id >> uint(i*8))
		}
	}
	l.size -= int64(n)
	return
}

func (l *FakeReader) WriteTo(w io.Writer) (n int64, err error) {
	size := int(l.size)
	bufferSize := len(l.sharedBytes)
	for size > 0 {
		tempBuffer := l.sharedBytes
		if size < bufferSize {
			tempBuffer = l.sharedBytes[0:size]
		}
		count, e := w.Write(tempBuffer)
		if e != nil {
			return int64(size), e
		}
		size -= count
	}
	return l.size, nil
}

func Readln(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return ln, err
}
