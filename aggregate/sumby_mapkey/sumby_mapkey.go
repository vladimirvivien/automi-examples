package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	// Declare the data source
	data := []map[string]string{
		{"Zip": "10452", "Count": "17", "Female": "12", "Male": "5"},
		{"Zip": "10453", "Count": "14", "Female": "7", "Male": "7"},
		{"Zip": "10454", "Count": "18", "Female": "8", "Male": "10"},
		{"Zip": "10455", "Count": "27", "Female": "17", "Male": "10"},
		{"Zip": "10456", "Count": "5", "Female": "3", "Male": "2"},
		{"Zip": "10458", "Count": "52", "Female": "25", "Male": "27"},
		{"Zip": "10459", "Count": "7", "Female": "5", "Male": "2"},
		{"Zip": "10460", "Count": "27", "Female": "20", "Male": "7"},
		{"Zip": "10461", "Count": "49", "Female": "26", "Male": "23"},
	}

	// creates a stream from the source
	strm := stream.From(sources.Slice(data))

	// run stream with operators
	strm.Run(
		// Reduce each incoming map to a summary map
		exec.Map(func(_ context.Context, data map[string]string) map[string]int {
			count, _ := strconv.Atoi(data["Count"])
			female, _ := strconv.Atoi(data["Female"])
			male, _ := strconv.Atoi(data["Male"])
			return map[string]int{"Count": count, "Female": female, "Male": male}
		}),

		// batch all data into a single window
		window.Batch[map[string]int](),

		// sums the "Female" key from the map
		exec.SumByMapKey[[]map[string]int]("Female"),
	)

	// Declare a sink to display result
	strm.Into(sinks.Func(func(items float64) error {
		fmt.Printf("Total female is %.0f\n", items)
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
