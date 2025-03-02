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
	type log map[string]string
	source := sources.Slice([]log{
		{"Event": "request", "Src": "/i/a", "Device": "00:11:51:AA", "Result": "accepted"},
		{"Event": "response", "Src": "/i/a/", "Device": "00:11:51:AA", "Result": "served"},
		{"Event": "request", "Src": "/i/b", "Device": "00:11:22:33", "Result": "accepted"},
		{"Event": "response", "Src": "/i/b", "Device": "00:11:22:33", "Result": "served"},
		{"Event": "request", "Src": "/i/c", "Device": "00:11:51:AA", "Result": "accepted"},
		{"Event": "response", "Src": "/i/c", "Device": "00:11:51:AA", "Result": "served"},
		{"Event": "request", "Src": "/i/d", "Device": "00:BB:22:DD", "Result": "accepted"},
		{"Event": "response", "Src": "/i/d", "Device": "00:BB:22:DD", "Result": "served"},
	})

	strm := stream.From(source)
	strm.WithLogSink(sinks.SlogJSON(slog.LevelDebug))
	strm.Run(
		// Filter response events
		exec.Filter(func(_ context.Context, e log) bool {
			return (e["Event"] == "response")
		}),

		// map type log -> map[string]string to support aggregate map operations downstream
		exec.Map(func(_ context.Context, e log) map[string]string {
			return e
		}),

		// Batch the data into --> []log
		window.Batch[map[string]string](),

		// Sort by the "Src" key
		exec.SortByMapKey[[]map[string]string]("Src"),
	)

	// Send stream items into a sink handler
	strm.Into(sinks.Func(func(sorted []map[string]string) error {
		for _, item := range sorted {
			fmt.Printf("%v\n", item)
		}
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
