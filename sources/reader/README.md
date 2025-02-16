## `io.Reader` Source

This example shows how to use Automi to stream data from a Go
`io.Reader` source. First, the code defines a new io.Reader with
string value as the source. Then it configures an Automi source
that will stream data as 30-byte chunks:

```go
data := `...`

// Setup data source as an io.Reader
reader := strings.NewReader(data)
// configure reader to stream 30-byte chunks
src := sources.Reader(reader).BufferSize(30)
```

Next, the code creates a stream from the source and defines a
`Map` and a `Filter` operator to process the data:

```go
// Create stream from source reading 30-byte chunks
strm := stream.From(src)

// Define stream operators to execute
strm.Run(
	exec.Map(func(_ context.Context, chunk []byte) string {
		str := string(chunk)
		return str
	}),

	// filter out requests
	exec.Filter(func(_ context.Context, e string) bool {
		return (strings.Contains(e, `"response"`))
	}),
)

// sink result in a collector function which prints it
strm.Into(sinks.Func(func(data string) error {
	fmt.Println(data)
	return nil
}))

```