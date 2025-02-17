package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/vladimirvivien/automi/collectors"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	data := `"request", "/i/a", "00:11:51:AA", "accepted"
	"response", "/i/a/", "00:11:51:AA", "served"
	"request", "/i/b", "00:11:22:33", "accepted"
	"response", "/i/b", "00:11:22:33", "served"
	"request", "/i/c", "00:11:51:AA", "accepted"
	"response", "/i/c", "00:11:51:AA", "served"
	"request", "/i/d",  "00:BB:22:DD", "accepted"
	"response", "/i/a", "00:BB:22:DD", "served"`

	// stream data line-by-line
	reader := strings.NewReader(data)
	stream := stream.New(sources.Scanner(reader, bufio.ScanLines))

	stream.Map(func(chunk interface{}) string {
		str := chunk.(string)
		return str
	})

	// filter out requests
	stream.Filter(func(e string) bool {
		return (strings.Contains(e, `"response"`))
	})

	// sink result in a collector function which prints it
	stream.Into(collectors.Func(func(data interface{}) error {
		e := data.(string)
		fmt.Println(e)
		return nil
	}))

	// open the stream
	if err := <-stream.Open(); err != nil {
		fmt.Println(err)
		return
	}
}
