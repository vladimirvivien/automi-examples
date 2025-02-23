package main

import (
	"context"
	"fmt"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

type log map[string]string

func main() {
	ch := make(chan log)
	go func() {
		ch <- log{"Event": "request", "path": "/i/a", "Device": "00:11:51:AA", "Result": "accepted"}
		ch <- log{"Event": "response", "path": "/i/a/", "Device": "00:11:51:AA", "Result": "served"}
		ch <- log{"Event": "request", "path": "/i/b", "Device": "00:11:22:33", "Result": "accepted"}
		ch <- log{"Event": "response", "path": "/i/b", "Device": "00:11:22:33", "Result": "served"}
		ch <- log{"Event": "request", "path": "/i/c", "Device": "00:11:51:AA", "Result": "accepted"}
		ch <- log{"Event": "response", "path": "/i/c", "Device": "00:11:51:AA", "Result": "served"}
		ch <- log{"Event": "request", "path": "/i/d", "Device": "00:BB:22:DD", "Result": "accepted"}
		ch <- log{"Event": "response", "path": "/i/a", "Device": "00:BB:22:DD", "Result": "served"}

		close(ch)
	}()

	// Create stream from a Go channel as source
	strm := stream.From(sources.Chan(ch))

	// Declare stream operations
	strm.Run(
		exec.Filter(func(_ context.Context, e log) bool {
			return (e["Event"] == "response")
		}),
	)

	// Define user-function to handle items
	strm.Into(sinks.Func(func(data log) error {
		fmt.Println(data)
		return nil
	}))

	// Open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
