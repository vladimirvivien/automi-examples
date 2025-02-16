package main

import (
	"context"
	"fmt"

	"github.com/vladimirvivien/automi/collectors"
	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	type log map[string]string
	src := sources.Slice([]log{
		{"Event": "request", "Src": "/i/a", "Device": "00:11:51:AA", "Result": "accepted"},
		{"Event": "response", "Src": "/i/a/", "Device": "00:11:51:AA", "Result": "served"},
		{"Event": "request", "Src": "/i/b", "Device": "00:11:22:33", "Result": "accepted"},
		{"Event": "response", "Src": "/i/b", "Device": "00:11:22:33", "Result": "served"},
		{"Event": "request", "Src": "/i/c", "Device": "00:11:51:AA", "Result": "accepted"},
		{"Event": "response", "Src": "/i/c", "Device": "00:11:51:AA", "Result": "served"},
		{"Event": "request", "Src": "/i/d", "Device": "00:BB:22:DD", "Result": "accepted"},
		{"Event": "response", "Src": "/i/d", "Device": "00:BB:22:DD", "Result": "served"},
	})

	stream := stream.From(src).Flow(
		exec.Filter(func(ctx context.Context, e log) bool {
			return (e["Event"] == "response")
		}),
		window.New[log](nil)
	)


	// GroupByKey returns a []map[group-key][]group-items
	stream.Batch().GroupByKey("Device")

	stream.Into(collectors.Func(func(data interface{}) error {
		items := data.([]map[interface{}][]interface{})
		for _, item := range items {
			for k, v := range item {
				fmt.Printf("%v = %v\n", k, v)
			}
		}
		return nil
	}))

	// open the stream
	if err := <-stream.Open(); err != nil {
		fmt.Println(err)
		return
	}
}
