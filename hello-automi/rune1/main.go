package main

import (
	"context"
	"fmt"
	"os"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {

	// create new stream with a slice of runes as source
	strm := stream.From(sources.Slice([]rune(`B世!ぽ@opqDQRS#$%^&*()ᅖ4x5Њ8yzUd90E12a3ᇳFGHmIザJuKLMᇙNO6PTnVWXѬYZbcef7ghijCklrAstvw`)))

	strm.Run(
		// filter out lowercase, non printable chars
		exec.Filter(func(_ context.Context, item rune) bool {
			return item >= 65 && item < (65+26)
		}),

		// map each rune to string value
		exec.Map(func(_ context.Context, item rune) string {
			return string(item)
		}),
	)

	// route string charaters to Stdout using a collector
	strm.Into(sinks.Writer[string](os.Stdout))

	// start the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
