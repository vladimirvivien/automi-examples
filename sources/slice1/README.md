## Slice Source (again)

In this Automi Slice source example, the code shows how to stream data from a Go slice of custom type items.
The code first creates an Automi source using a slice of `log` values as its data source. Then the remainde of the code applies stream operators to process the data:

```go
type log map[string]string

data := []log{...}

// Declare stream with a Slice source
strm := stream.From(sources.Slice(data))

// Define stream operations
strm.Run(
	exec.Filter(func(_ context.Context, e log) bool {
		return (e["Event"] == "response")
	}),
)

// Send streamed item to a sink
strm.Into(sinks.Func(func(data log) error {
	fmt.Println(data)
	return nil
}))
```