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
	// Setup data source as slice of slices
	source := sources.Slice([][]string{
		{"request", "/i/a", "00:11:51:AA", "accepted"},
		{"response", "/i/a/", "00:11:51:AA", "failed"},
		{"request", "/i/b", "00:11:22:33", "accepted"},
		{"response", "/i/b", "00:11:22:33", "served"},
		{"request", "/i/c", "00:11:51:AA", "accepted"},
		{"response", "/i/c", "00:11:51:AA", "served"},
		{"request", "/i/d", "00:BB:22:DD", "accepted"},
		{"response", "/i/d", "00:BB:22:DD", "failed"},
	})

	// Create new stream from source
	strm := stream.From(source)

	// Setup stream operations
	strm.Run(
		// Batch each incoming `[]string` into --> [][]string
		window.Batch[[]string](),

		// Group batch items by value at index 3 into --> map[string][][]string
		exec.GroupByIndex[[][]string](3),
	)

	// Define sink operator to handle items
	strm.Into(sinks.Func(func(items map[string][][]string) error {
		for group, item := range items {
			fmt.Println("Group:", group)
			for _, slice := range item {
				fmt.Printf("%v\n", slice)
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
