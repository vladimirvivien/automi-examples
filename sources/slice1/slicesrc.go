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

	data := []log{
		log{"Event": "request", "path": "/i/a", "Device": "00:11:51:AA", "Result": "accepted"},
		log{"Event": "response", "path": "/i/a/", "Device": "00:11:51:AA", "Result": "served"},
		log{"Event": "request", "path": "/i/b", "Device": "00:11:22:33", "Result": "accepted"},
		log{"Event": "response", "path": "/i/b", "Device": "00:11:22:33", "Result": "served"},
		log{"Event": "request", "path": "/i/c", "Device": "00:11:51:AA", "Result": "accepted"},
		log{"Event": "response", "path": "/i/c", "Device": "00:11:51:AA", "Result": "served"},
		log{"Event": "request", "path": "/i/d", "Device": "00:BB:22:DD", "Result": "accepted"},
		log{"Event": "response", "path": "/i/a", "Device": "00:BB:22:DD", "Result": "served"},
	}

	// Declare stream with a Slice source
	strm := stream.From(sources.Slice(data))

	// Define stream operations
	strm.Run(
		exec.Filter(func(_ context.Context, e log) bool {
			return (e["Event"] == "response")
		}),
	)

	// Send streamed item to a sink
	strm.Into(sinks.Func(func(data log) error {
		fmt.Println(data)
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
