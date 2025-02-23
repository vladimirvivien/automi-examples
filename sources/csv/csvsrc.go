package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

// Loads CSV data from file.
func main() {
	srcFile, err := os.Open("./stats.csv")
	if err != nil {
		slog.Error("Unable to open file", "error", err)
		os.Exit(1)
	}
	defer srcFile.Close()

	// Configure the CSV source
	csv := sources.CSV(srcFile).
		CommentChar('#'). // sets comment charcter, default is '#'
		DelimChar(',').   // sets delimiter char, default is  ','
		HasHeaders()      // indicate CSV has headers

	// Create new stream from CSV source
	strm := stream.From(csv)

	strm.Run(
		// select first 6 cols per row
		exec.Map(func(_ context.Context, row []string) []string {
			return row[:6]
		}),

		// filter out rows with zero participants
		exec.Filter(func(_ context.Context, row []string) bool {
			count, err := strconv.Atoi(row[1])
			if err != nil {
				count = 0
			}
			return (count > 0)
		}),
	)

	// sink result in a collector function to prints it
	strm.Into(sinks.Func(func(row []string) error {
		fmt.Println(row)
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		slog.Error("Streamed failed", "error", err)
		os.Exit(1)
	}
}
