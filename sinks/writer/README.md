## io.Writer Sink

This example shows how to define an Automi stream
to forward streamed items into a sink that stores
the items using an io.Writer. 

First define the slink with the slice as a storage:

```go
sink := new(bytes.Buffer)
```

Then define a stream that uses a Go slice as the sink:

```go
// Create a stream from Go channel source
strm := stream.From(sources.Chan(ch))

// Define the stream operations
strm.Run(
	exec.Map(func(_ context.Context, row string) []string {
		return strings.Split(row, ",")
	}),
)

// Set up a stream sink with the buffer
strm.Into(sinks.Writer[[]byte](sink))
```

In the previous example, the `sinks.Writer` function
can be instantiated with either a `[]byte` or a `string`
type.

Lastly, the writer can be accessed:

```go
fmt.Println(sink.String())
```