package log

import (
	"bytes"
	"io"
	"sync"
)

type BufferWriter struct {
	MaxBufferSize int64
	sync.Mutex
	Buffer *bytes.Buffer
	Writer io.Writer
}

func NewBuffer(writer io.Writer, bufferSize int64) *BufferWriter {
	if this, ok := writer.(*BufferWriter); ok {
		return this
	}
	return &BufferWriter{
		Writer:        writer,
		MaxBufferSize: bufferSize,
		Buffer:        bytes.NewBuffer(make([]byte, 0, bufferSize)),
	}
}

func (this *BufferWriter) Write(p []byte) (n int, err error) {
	this.Lock()
	defer this.Unlock()
	p = append(p, '\n')
	needBufferSize := len(p)
	if int64(needBufferSize) >= this.MaxBufferSize {
		err := this.bufferToWriter()
		if err != nil {
			return 0, err
		}
		return this.Writer.Write(p)
	}
	needBufferSize = this.Buffer.Len() + len(p)
	if int64(needBufferSize) >= this.MaxBufferSize {
		err := this.bufferToWriter()
		if err != nil {
			return 0, err
		}
	}
	return this.Buffer.Write(p)
}

func (this *BufferWriter) bufferToWriter() error {
	_, err := this.Buffer.WriteTo(this.Writer)
	return err
}
