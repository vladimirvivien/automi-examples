## Go `slog` Sink

This example shows how to use the `slog` Logger as
a sink for items in an Automi stream. The `slog` sink
is designed to be used as a side stream to receive
log events (see example [here](../../logging/)). However,
this example shows how it can be used as a general sink.

First, let's define a sink backed by a `slog.Handler`:

```go
logger := sinks.SlogJSON(slog.LevelDebug)
```

Next, define the stream. Note that the `exec.Execute` operation
prepares the stream item as a `api.StreamLog` type to be 
consumed by the `slog` sink:

```go
data := []string{...}

// Define the stream with a Go slice source
strm := stream.From(sources.Slice(data))

// Define the stream operations
strm.Run(
	exec.Execute(func(ctx context.Context, row string) api.StreamLog {
		data := strings.Split(row, " ")
		return api.StreamLog{
			Message: "Item execution",
			Level:   slog.LevelDebug,
			Attrs:   []slog.Attr{slog.Any("data", data)},
		}
	}),
)

// Setup the slice as stream sink
strm.Into(logger)
```