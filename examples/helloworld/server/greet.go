package main

import (
	"context"
	"net/http"
	"time"

	"github.com/weiwenchen2022/wsprotobuf"
	"nhooyr.io/websocket"

	pb "github.com/weiwenchen2022/wsprotobuf/examples/helloworld/helloworld"
)

// greetServer is the WebSocket greet server implementation.
type greetServer struct {
	// logf controls where logs are sent.
	logf func(string, ...any)
}

func (s greetServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.CloseNow()

	for {
		err = s.greet(r.Context(), c)
		switch websocket.CloseStatus(err) {
		case websocket.StatusNormalClosure, websocket.StatusGoingAway:
			return
		}
		if err != nil {
			s.logf("failed to echo with %s: %v", r.RemoteAddr, err)
			return
		}
	}
}

// greet reads from the WebSocket connection and then writes
// the response message back to it.
// The entire function has 10s to complete.
func (s greetServer) greet(ctx context.Context, c *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var in pb.HelloRequest
	err := wsprotobuf.Read(ctx, c, &in)
	if err != nil {
		return err
	}

	s.logf("Received: %s", in.GetName())
	return wsprotobuf.Write(ctx, c, &pb.HelloReply{Message: "Hello " + in.GetName()})
}
