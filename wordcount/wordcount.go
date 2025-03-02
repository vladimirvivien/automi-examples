package main

import (
	"bufio"
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/vladimirvivien/automi/api/tuple"
	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/flat"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	regSpaces := regexp.MustCompile(`\s+`)
	regNonAlpha := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	space := " "

	// setup stream data source from file
	file, err := os.Open("./twotw.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	strm := stream.From(sources.Scanner(file, bufio.ScanLines))
	strm.WithLogSink(sinks.SlogJSON(slog.LevelDebug))
	strm.Run(
		// normalize all spaces to single space, then return as slice of words
		exec.Execute(func(_ context.Context, line []byte) []string {
			return strings.Split(regSpaces.ReplaceAllLiteralString(string(line), space), space)
		}),

		// Flatten the slice of words into individual words
		flat.Slice[[]string](),

		// Map each word to an occurence tuple of Pair{word, 1}
		exec.Map(func(_ context.Context, word string) tuple.Pair[string, int] {
			word = regNonAlpha.ReplaceAllLiteralString(word, "")
			return tuple.Pair[string, int]{Val1: word, Val2: 1}
		}),

		// Batch the occurence pairs into slices of []Pair{word, 1}
		window.Batch[tuple.Pair[string, int]](),

		// Group the occurences by word --> map[word][]Pair{word, 1}:
		exec.GroupByStructField[[]tuple.Pair[string, int]]("Val1"),

		// Count the occurences of each word from the group, return as []Pair{word, total}
		exec.Execute(func(_ context.Context, group map[any][]tuple.Pair[string, int]) []tuple.Pair[string, int] {
			var result []tuple.Pair[string, int]
			for key, pairs := range group {
				sum := 0
				for _, pair := range pairs {
					sum += pair.Val2
				}
				result = append(result, tuple.Pair[string, int]{Val1: key.(string), Val2: sum})
			}

			// sort result
			slices.SortFunc(result, func(p1, p2 tuple.Pair[string, int]) int {
				return cmp.Compare(p1.Val1, p2.Val1)
			})

			return result
		}),
	)

	// Route the result to a sink that prints the word count
	strm.Into(sinks.Func(func(wordp []tuple.Pair[string, int]) error {
		for _, pair := range wordp {
			fmt.Printf("%s: %d\n", pair.Val1, pair.Val2)
		}
		return nil
	}))

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
