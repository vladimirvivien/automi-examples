package main

import (
	"bufio"
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

	// Define io.Scanner source to read data line-by-line
	reader := strings.NewReader(data)
	src := sources.Scanner(reader, bufio.ScanLines)

	// Create stream from io.Scanner
	strm := stream.From(src)

	// Define stream operations
	strm.Run(
		// map line to string
		exec.Map(func(_ context.Context, chunk []byte) string {
			return string(chunk)
		}),

		// filter out requests
		exec.Filter(func(_ context.Context, e string) bool {
			return (strings.Contains(e, `"response"`))
		}),
	)

	// sink result in a collector function which prints it
	strm.Into(sinks.Func(func(data string) error {
		fmt.Println(strings.TrimSpace(data))
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
