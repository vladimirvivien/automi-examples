package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

// scientist represents row data in source (csv file)
type scientist struct {
	FirstName string
	LastName  string
	Title     string
	BornYear  int
}

// Example:
// - load data from a csv
// - map row to a custom type
// - filter based on selected value
// - write out result to a file
func main() {
	// prepare data source
	src, err := os.Open("./data.txt")
	if err != nil {
		fmt.Println("Unable to open source:", err)
		os.Exit(1)
	}
	defer src.Close()
	source := sources.CSV(src)

	// prepare data
	snk, err := os.Create("./result.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer snk.Close()
	sink := sinks.CSV(snk)

	// start stream definition
	stream := stream.From(source)

	stream.Run(
		// map csv row to struct scientist
		exec.Map(func(ctx context.Context, cs []string) scientist {
			yr, _ := strconv.Atoi(cs[3])
			return scientist{
				FirstName: cs[1],
				LastName:  cs[0],
				Title:     cs[2],
				BornYear:  yr,
			}
		}),

		// apply data filter
		exec.Filter(func(ctx context.Context, cs scientist) bool {
			return (cs.BornYear > 1930)
		}),

		// remap value of type scientst to []string
		exec.Map(func(ctx context.Context, cs scientist) []string {
			return []string{cs.FirstName, cs.LastName, cs.Title}
		}),
	)
	// stream result into sink
	stream.Into(sink)

	if err := <-stream.Open(context.TODO()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("wrote result to file result.txt")
}
