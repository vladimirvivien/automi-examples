package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	// Create a channel source
	ch := make(chan string)
	go func() {
		defer close(ch)
		ch <- "10452,17,12,0.71,5,0.29,0,0,17,100"
		ch <- "10453,14,7,0.5,7,0.5,0,0,14,100"
		ch <- "10454,18,8,0.44,10,0.56,0,0,18,100"
		ch <- "10455,27,17,0.63,10,0.37,0,0,27,100"
		ch <- "10456,5,3,0.6,2,0.4,0,0,5,100"
		ch <- "10458,52,25,0.48,27,0.52,0,0,52,100"
		ch <- "10459,7,5,0.71,2,0.29,0,0,7,100"
		ch <- "10460,27,20,0.74,7,0.26,0,0,27,100"
		ch <- "10461,49,26,0.53,23,0.47,0,0,49,100"
	}()

	// Create a stream from channel source
	strm := stream.From(sources.Chan(ch))

	// Define stream operations
	strm.Run(
		// Split each row into fields
		exec.Map(func(_ context.Context, row string) []string {
			return strings.Split(row, ",")
		}),

		// Extract the 4th field and convert to float
		exec.Map(func(_ context.Context, data []string) float64 {
			f, _ := strconv.ParseFloat(data[3], 32)
			return f
		}),

		// Batch incoming float64 --> []float64
		window.Batch[float64](),

		// Sum all float64 values into --> a single float64
		exec.Sum[[]float64, float64](),
	)

	// Define sink to print total
	strm.Into(sinks.Func(func(total float64) error {
		fmt.Printf("Total is %.2f\n", total)
		return nil
	}))

	// Open stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
