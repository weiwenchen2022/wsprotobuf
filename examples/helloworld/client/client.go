package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/weiwenchen2022/wsprotobuf"
	"nhooyr.io/websocket"

	pb "github.com/weiwenchen2022/wsprotobuf/examples/helloworld/helloworld"
)

const defaultName = "world"

var name = flag.String("name", defaultName, "Name to greet")

func main() {
	flag.Parse()
	log.SetFlags(0)

	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s URL", os.Args[0])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	c, _, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		log.Fatalf("failed to dial with %s: %v", u.String(), err)
	}
	defer c.CloseNow()

	err = wsprotobuf.Write(ctx, c, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("failed to write: %v", err)
	}

	var r pb.HelloReply
	err = wsprotobuf.Read(ctx, c, &r)
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())

	c.Close(websocket.StatusNormalClosure, "")
}
