// Package wsprotobuf provides helpers for reading and writing protobuf messages.
package wsprotobuf

import (
	"context"
	"fmt"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

// Read reads a protobuf message from c into v.
// It will reuse buffers in between calls to avoid allocations.
func Read(ctx context.Context, c *websocket.Conn, v any) (err error) {
	m, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("%T not implements proto.Message interface", v)
	}

	_, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	b, put := getBuffer()
	defer put()

	_, err = io.Copy(b, r)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(b.bytes(), m)
	if err != nil {
		c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal protobuf")
		return err
	}
	return nil
}

// Write writes the protobuf message v to c.
// It will reuse buffers in between calls to avoid allocations.
func Write(ctx context.Context, c *websocket.Conn, v any) error {
	m, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("%T not implements proto.Message interface", v)
	}

	b, put := getBuffer()
	defer put()
	nb, err := (proto.MarshalOptions{}).MarshalAppend(*(*[]byte)(b), m)
	*b = buf(nb)
	if err != nil {
		return err
	}
	return c.Write(ctx, websocket.MessageBinary, b.bytes())
}

var bufPool = sync.Pool{
	New: func() any {
		return new(buf)
	},
}

func getBuffer() (b *buf, put func()) {
	b = bufPool.Get().(*buf)
	b.reset()
	return b, func() { bufPool.Put(b) }
}

type buf []byte

func (b *buf) reset() {
	*b = (*b)[:0]
}

func (b *buf) Write(p []byte) (int, error) {
	*b = buf(append(*(*[]byte)(b), p...))
	return len(p), nil
}

func (b *buf) bytes() []byte {
	return *(*[]byte)(b)
}
