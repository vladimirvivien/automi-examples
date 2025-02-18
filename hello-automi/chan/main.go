package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {

	// emitterFunc returns a chan used for data
	emitterFunc := func() <-chan time.Time {
		times := make(chan time.Time)
		go func() {
			times <- time.Unix(100000, 0)
			times <- time.Unix(2*100000, 0)
			times <- time.Unix(4*100000, 0)
			times <- time.Unix(8*100000, 0)
			close(times)
		}()
		return times
	}

	// Create stream from channel source
	strm := stream.From(sources.Chan(emitterFunc()))

	// Setup execution of operators
	strm.Run(
		// map each rune to string value
		exec.Map(func(_ context.Context, item time.Time) string {
			return item.String()
		}),
	)

	// route string charaters to Stdout using a collector
	strm.Into(sinks.Writer[string](os.Stdout))

	// start stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
