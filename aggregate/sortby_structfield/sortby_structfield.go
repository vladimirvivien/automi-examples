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
	type log struct{ Event, Src, Device, Result string }
	src := sources.Slice([]log{
		{Event: "request", Src: "/i/a", Device: "00:11:51:AA", Result: "accepted"},
		{Event: "response", Src: "/i/a/", Device: "00:11:51:AA", Result: "served"},
		{Event: "request", Src: "/i/b", Device: "00:11:22:33", Result: "accepted"},
		{Event: "response", Src: "/i/b", Device: "00:11:22:33", Result: "served"},
		{Event: "request", Src: "/i/c", Device: "00:11:51:AA", Result: "accepted"},
		{Event: "response", Src: "/i/c", Device: "00:11:51:AA", Result: "served"},
		{Event: "request", Src: "/i/d", Device: "00:BB:22:DD", Result: "accepted"},
		{Event: "response", Src: "/i/d", Device: "00:BB:22:DD", Result: "served"},
	})

	// Define new stream with operations chained
	strm := stream.From(src).Run(
		
		// Filter out just responses
		exec.Filter(func(ctx context.Context, e log) bool {
			return (e.Event == "response")
		}),

		// Batch logs into -> []log
		window.Batch[log](),
		
		// Sort batched items by Src field
		exec.SortByStructField[[]log]("Src"),
	)

	// Set up sink to process sorted results
	strm.Into(sinks.Func(func(items []log) error {
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
