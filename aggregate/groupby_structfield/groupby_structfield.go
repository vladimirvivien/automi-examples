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

type log struct{ Event, Src, Device, Result string }

func main() {
	// Declare data source
	source := sources.Slice([]log{
		{Event: "request", Src: "/i/a", Device: "00:11:51:AA", Result: "accepted"},
		{Event: "response", Src: "/i/a", Device: "00:11:51:AA", Result: "served"},
		{Event: "request", Src: "/i/b", Device: "00:11:22:33", Result: "accepted"},
		{Event: "response", Src: "/i/b", Device: "00:11:22:33", Result: "served"},
		{Event: "request", Src: "/i/c", Device: "00:11:51:AA", Result: "accepted"},
		{Event: "response", Src: "/i/c", Device: "00:11:51:AA", Result: "served"},
		{Event: "request", Src: "/i/d", Device: "00:BB:22:DD", Result: "accepted"},
		{Event: "response", Src: "/i/d", Device: "00:BB:22:DD", Result: "served"},
		{Event: "response", Src: "/i/a", Device: "00:12:51:AA", Result: "accepted"},
	})

	// Create new stream from source
	strm := stream.From(source)

	// Setup stream operations
	strm.Run(
		// Batch each incoming `log` items into --> `[]log`
		window.Batch[log](),

		// Group incoming `[]log` by struct field `Src` into --> map[any][]log
		exec.GroupByStructField[[]log]("Src"),
	)

	// Define sink opertor to handle items
	strm.Into(sinks.Func(func(items map[any][]log) error {
		for key, item := range items {
			fmt.Println("Group:", key)
			fmt.Println("-----")
			for _, log := range item {
				fmt.Printf("log = %v\n", log)
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
