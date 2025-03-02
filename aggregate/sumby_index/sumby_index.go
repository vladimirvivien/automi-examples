package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	// Open data source file
	file, err := os.Open("stats.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Define an Automi CSV source from file
	src := sources.CSV(file).
		CommentChar('#'). // sets comment charcter, default is '#'
		DelimChar(',').   // sets delimiter char, default is  ','
		HasHeaders()      // indicate CSV has headers

	// Create a stream from source
	strm := stream.From(src)

	// Define stream operations
	strm.Run(
		// Reduces and maps each row to []int containing summary data
		exec.Map(func(_ context.Context, data []string) []int {
			count, _ := strconv.Atoi(data[1])
			female, _ := strconv.Atoi(data[2])
			male, _ := strconv.Atoi(data[4])
			return []int{count, female, male}
		}),

		// Batch incoming []int into a window
		window.Batch[[]int](),

		// Sum all values at pos(1)
		exec.SumByIndex[[][]int](1),
	)

	strm.Into(sinks.Func(func(val float64) error {
		fmt.Printf("Total female is %.0f\n", val)
		return nil
	}))

	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
