## Slice Sink

This example shows how to define an Automi stream
to forward streamed items into a sink that stores
the items into a slice. Items from the sink slice 
can then be accessed.

First define the slink with the slice as a storage:

```go
slice := sinks.Slice[[]string]()
```

Then define a stream that uses a Go slice as the sink:

```go
data := []string{...}

// Define the stream with a Go slice source
strm := stream.From(sources.Slice(data))

// Define the stream operations
strm.Run(
	exec.Map(func(_ context.Context, row string) []string {
		return strings.Split(row, " ")
	}),
)

// Setup the slice as stream sink
strm.Into(slice)
```

The data in the sink slice can be accessed as follows:

```go
for _, item := range slice.Get() {
	fmt.Println(item)
}
```