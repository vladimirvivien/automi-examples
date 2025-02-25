package main

import (
	"cmp"
	"context"
	"fmt"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

// Define stream item type
type log struct{ Event, Src, Device, Result string }

func main() {
	// Setup a data source for stream
	source := sources.Slice([]log{
		{Event: "request", Src: "/i/a", Device: "00:11:51:AA", Result: "accepted"},
		{Event: "response", Src: "/i/a/", Device: "00:11:51:AA", Result: "served"},
		{Event: "request", Src: "/i/b", Device: "00:11:22:33", Result: "accepted"},
		{Event: "response", Src: "/i/b", Device: "00:11:22:33", Result: "served"},
		{Event: "request", Src: "/i/c", Device: "00:11:51:AA", Result: "accepted"},
		{Event: "response", Src: "/i/c", Device: "00:11:51:AA", Result: "served"},
		{Event: "request", Src: "/i/d", Device: "00:BB:22:DD", Result: "accepted"},
		{Event: "response", Src: "/i/d", Device: "00:BB:22:DD", Result: "served"},
	})

	// Create new stream using source
	strm := stream.From(source)

	// Define stream operators
	strm.Run(
		// Filter some data
		exec.Filter(func(_ context.Context, e log) bool {
			return (e.Event == "response")
		}),

		// Batch the data into --> []log
		window.Batch[log](),

		// Define a sort function for the streamed items
		exec.SortWithFunc[[]log](func(i, j log) int {
			return cmp.Compare(i.Src, j.Src)
		}),
	)

	// Send stream items into a sink handler
	strm.Into(sinks.Func(func(sorted []log) error {
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
