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

// Streams data from CSV source and stores collected result in CSV
func main() {
	// setup source and sink files
	srcFile, err := os.Open("./stats.txt")
	if err != nil {
		slog.Error("Failed to open file", "err", err)
		os.Exit(1)
	}
	defer srcFile.Close()

	sinkFile, err := os.Create("./result.txt")
	if err != nil {
		slog.Error("Failed to create result file", "err", err)
		os.Exit(1)
	}
	defer sinkFile.Close()

	source := sources.CSV(srcFile).
		CommentChar('#').
		DelimChar(',').
		HasHeaders()

	// create stream from CSV file
	strm := stream.From(source)

	// Run stream operators
	strm.Run(
		// select first 6 cols per row:
		exec.Map(func(_ context.Context, row []string) []string {
			return row[:6]
		}),
		// filter out rows with col[1] = 0
		exec.Filter(func(_ context.Context, row []string) bool {
			count, err := strconv.Atoi(row[1])
			if err != nil {
				count = 0
			}
			return (count > 0)
		}),
	)
	strm.Into(sinks.CSV(sinkFile))

	// launch the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
