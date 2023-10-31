package main

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/weiwenchen2022/wsprotobuf"
	pb "github.com/weiwenchen2022/wsprotobuf/examples/helloworld/helloworld"
	"nhooyr.io/websocket"
)

// Test_greetServer tests the greetServer by sending it 5 different messages
// and ensuring the responses all match.
func Test_greetServer(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(greetServer{
		logf: t.Logf,
	})
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for i := 0; i < 5; i++ {
		err = wsprotobuf.Write(ctx, c, &pb.HelloRequest{Name: strconv.Itoa(i)})
		if err != nil {
			t.Fatal(err)
		}

		var r pb.HelloReply
		err = wsprotobuf.Read(ctx, c, &r)
		if err != nil {
			t.Fatal(err)
		}

		got := r.GetMessage()
		expected := fmt.Sprintf("Hello %d", i)

		if expected != got {
			t.Fatalf("expected %q but got %q", expected, got)
		}
	}

	c.Close(websocket.StatusNormalClosure, "")
}
