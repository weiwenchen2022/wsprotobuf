// Package wsprotobuf provides helpers for reading and writing protobuf messages.
package wsprotobuf

import (
	"bytes"
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

	buf, put := getBuffer()
	defer put()

	_, err = io.Copy(buf, r)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(buf.Bytes(), m)
	if err != nil {
		c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal protobuf")
		return err
	}
	return nil
}

// Write writes the protobuf message v to c.
func Write(ctx context.Context, c *websocket.Conn, v any) error {
	m, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("%T not implements proto.Message interface", v)
	}

	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return c.Write(ctx, websocket.MessageBinary, b)
}

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func getBuffer() (buf *bytes.Buffer, put func()) {
	buf = bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf, func() { bufPool.Put(buf) }

}
