## User-defined Function Sink

This example shows how to define an Automi stream
to forward streamed items into sink to be processed
using a user-defined function.

The example sources it's data from a Go channel,
then maps each streamed item to a string slice.
The slice is forwarded to the sink where it is
processed by the user-provided function.

```go
// Define a new stream from a Go channel source
strm := stream.From(sources.Chan(ch))

// Define stream opeations to run
strm.Run(
	exec.Map(func(_ context.Context, row string) []string {
		return strings.Split(row, ",")
	}),
)

// Define stream sink as a user-provided function
strm.Into(sinks.Func(func(row []string) error {
	fmt.Println(row[len(row)-1])
	return nil
}))
```
