package main

import (
	"context"
	"fmt"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	src := sources.Slice([][]string{
		{"request", "/i/a", "00:11:51:AA", "accepted"},
		{"response", "/i/a/", "00:11:51:AA", "failed"},
		{"request", "/i/b", "00:11:22:33", "accepted"},
		{"response", "/i/b", "00:11:22:33", "served"},
		{"request", "/i/c", "00:11:51:AA", "accepted"},
		{"response", "/i/c", "00:11:51:AA", "served"},
		{"request", "/i/d", "00:BB:22:DD", "accepted"},
		{"response", "/i/d", "00:BB:22:DD", "failed"},
	})

	// Define new stream with operations chained
	strm := stream.From(src).Run(
		// Filter out just responses
		exec.Filter(func(ctx context.Context, e []string) bool {
			return (e[0] == "response")
		}),

		// Batch items into -> [][]string
		window.Batch[[]string](),

		// Sort batched items by position index 3 (result field)
		exec.SortSliceByIndex[[][]string](3),
	)

	// Set up sink to process sorted results
	strm.Into(sinks.Func(func(items [][]string) error {
		for _, item := range items {
			fmt.Printf("%v\n", item)
		}
		return nil
	}))

	// Open the stream with context
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
