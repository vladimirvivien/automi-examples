## `bufio.Scanner` Source

This example shows how to use Automi to stream data from a Go
`bufio.Scanner` source. First, the code defines a new `io.Reader` with
a string value as its source. Then, it configures an Automi `Scanner` source
to read the data using a Go `bufio.Scanner` line by line:

```go
data := `...`

// Define io.Scanner source to read data line-by-line
reader := strings.NewReader(data)
src := sources.Scanner(reader, bufio.ScanLines)
```

Next, the code creates a stream from the source and defines a
`Map` and a `Filter` operator to process the data:

```go
// Create stream from io.Scanner
strm := stream.From(src)

// Define stream operations
strm.Run(
	// map line to string
	exec.Map(func(_ context.Context, chunk []byte) string {
		return string(chunk)
	}),

	// filter out requests
	exec.Filter(func(_ context.Context, e string) bool {
		return (strings.Contains(e, `"response"`))
	}),
)

// sink result in a collector function which prints it
strm.Into(sinks.Func(func(data string) error {
	fmt.Println(strings.TrimSpace(data))
	return nil
}))

```