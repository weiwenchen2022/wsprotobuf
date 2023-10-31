package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "localhost:0", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run starts a http.Server for the passed in address
// with all requests handled by greetServer.
func run() error {
	l, err := net.Listen("tcp", *addr)
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	s := &http.Server{
		Handler: greetServer{
			logf: log.Printf,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case <-sigc:
		log.Println("terminating")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}
