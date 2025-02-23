package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	src := sources.Slice([]map[string]string{
		{"Event": "request", "Src": "/i/a", "Device": "00:11:51:AA", "Result": "accepted"},
		{"Event": "response", "Src": "/i/a/", "Device": "00:11:51:AA", "Result": "served"},
		{"Event": "request", "Src": "/i/b", "Device": "00:11:22:33", "Result": "accepted"},
		{"Event": "response", "Src": "/i/b", "Device": "00:11:22:33", "Result": "served"},
		{"Event": "request", "Src": "/i/c", "Device": "00:11:51:AA", "Result": "accepted"},
		{"Event": "response", "Src": "/i/c", "Device": "00:11:51:AA", "Result": "served"},
		{"Event": "request", "Src": "/i/d", "Device": "00:BB:22:DD", "Result": "accepted"},
		{"Event": "response", "Src": "/i/d", "Device": "00:BB:22:DD", "Result": "served"},
	})

	// Define new stream from slice source
	strm := stream.From(src)

	// attach a logger sink
	strm.WithLogSink(sinks.SlogText(slog.LevelDebug))

	// define stream operations
	strm.Run(
		// filter out "response" event
		exec.Filter(func(ctx context.Context, e map[string]string) bool {
			return (e["Event"] == "response")
		}),

		// Batch map items into -> []map[string]string
		window.Batch[map[string]string](),

		// Group batched items []map[string]string -> map[string][]map[string]string
		exec.GroupByMapKey[[]map[string]string]("Device"),
	)

	strm.Into(sinks.Func(func(items map[string][]map[string]string) error {
		for batch, items := range items {
			fmt.Println("Batch:", batch)
			fmt.Println("-----")
			for _, item := range items {
				fmt.Printf("%#v\n", item)
			}
		}
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
