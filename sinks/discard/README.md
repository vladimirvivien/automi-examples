## The Discard Sink

As its name implies, this sink discards any items it
receives. It is useful when testing stream with large data.

The example starts by defining a stream from a Go channels source.
Each streamed item is mapped to a string slice then is forwarded
to the discard sink.

```go
// Define a new stream from a Go channel source
strm := stream.From(sources.Chan(ch))

// Define stream opeations to run
strm.Run(
	exec.Map(func(_ context.Context, row string) []string {
		return strings.Split(row, ",")
	}),
)

// Define a discard sink which is effectively a noop
strm.Into(sinks.Discard())
```