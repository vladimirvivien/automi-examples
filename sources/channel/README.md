## Channel Source

This example shows how to stream data from an Automi
source backed by a Go channel.

```go
ch := make(chan log)
... 

// Create stream from a Go channel as source
strm := stream.From(sources.Chan(ch))

// Declare stream operations
strm.Run(
	exec.Filter(func(_ context.Context, e log) bool {
		return (e["Event"] == "response")
	}),
)

// Define user-function to handle items
strm.Into(sinks.Func(func(data log) error {
	fmt.Println(data)
	return nil
}))
```