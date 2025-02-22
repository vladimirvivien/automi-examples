## CSV Sink

This example shows how to define an Automi stream
to forward streamed items into a CSV sink.

First, the CSV sink is defined to use an in-memory
storage (`bytes.Buffer`):

```go
var strBuilder = bytes.NewBufferString("")
csvSink := sinks.CSV(strBuilder)
csvSink.DelimChar('|')
```

Next, the example defines a simple stream to
split the incoming data into string slices to
be sent to the CSV sink:

```go
// Create a stream from Slice source
strm := stream.From(sources.Slice(data))

// Define stream operators to run
strm.Run(
	exec.Map(func(_ context.Context, row string) []string {
		return strings.Split(row, " ")
	}),
)

// Define stream sink
strm.Into(csvSink)
```