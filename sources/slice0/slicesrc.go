package main

import (
	"context"
	"fmt"
	"os"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/operators/window"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	// Define slice source
	slice := sources.Slice([]rune(`B世!ぽ@opqDQRS#$%^&*()ᅖ4x5Њ8yzUd90E12a3ᇳFGHmIザJuKLMᇙNO6PTnVWXѬYZbcef7ghijCklrAstvw`))

	// create stream with emitter of rune slice
	strm := stream.From(slice)

	strm.Run(
		exec.Filter(func(_ context.Context, item rune) bool {
			return item >= 65 && item < (65+26) // remove unwanted chars
		}),
		exec.Map(func(_ context.Context, item rune) string {
			return string(item) // map rune to string
		}),

		// batch incoming string items
		window.Batch[string](),

		// sort batched items
		exec.SortSlice[[]string](),
	)

	// Send result to stdout
	strm.Into(sinks.Writer[string](os.Stdout)) 

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
