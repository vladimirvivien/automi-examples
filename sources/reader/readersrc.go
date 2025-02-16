package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	data := `"request", "/i/a", "00:11:51:AA", "accepted"
	"response", "/i/a/", "00:11:51:AA", "served"
	"request", "/i/b", "00:11:22:33", "accepted"
	"response", "/i/b", "00:11:22:33", "served"
	"request", "/i/c", "00:11:51:AA", "accepted"
	"response", "/i/c", "00:11:51:AA", "served"
	"request", "/i/d",  "00:BB:22:DD", "accepted"
	"response", "/i/a", "00:BB:22:DD", "served"`

	// Setup data source as an io.Reader
	reader := strings.NewReader(data)
	// configure reader to stream 30-byte chunks
	src := sources.Reader(reader).BufferSize(30)

	// Create stream from source reading 30-byte chunks
	strm := stream.From(src)

	// Define stream operators to execute
	strm.Run(
		exec.Map(func(_ context.Context, chunk []byte) string {
			str := string(chunk)
			return str
		}),

		// filter out requests
		exec.Filter(func(_ context.Context, e string) bool {
			return (strings.Contains(e, `"response"`))
		}),
	)

	// sink result in a collector function which prints it
	strm.Into(sinks.Func(func(data string) error {
		fmt.Println(data)
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
