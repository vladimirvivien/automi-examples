Collecting workspace information# SortByStructField Example

This example demonstrates how to sort a batch of structured data by a specific struct field using Automi. The code creates a stream from a slice of custom log entries, filters for specific entries, batches them, and sorts by a specified field.

## Code walkthrough

First, the code defines a custom `log` struct type and creates a slice source:

```go
type log struct{ Event, Src, Device, Result string }
src := sources.Slice([]log{
    {Event: "request", Src: "/i/a", Device: "00:11:51:AA", Result: "accepted"},
    {Event: "response", Src: "/i/a/", Device: "00:11:51:AA", Result: "served"},
    // more log entries...
})
```

Next, it creates a stream from this source and chains operations using the `Run()` method:

```go
strm := stream.From(src).Run(
    // Filter out just responses
    exec.Filter(func(ctx context.Context, e log) bool {
        return (e.Event == "response")
    }),

    // Batch logs into -> []log
    window.Batch[log](),
    
    // Sort batched items by Src field
    exec.SortByStructField[[]log]("Src"),
)
```

The stream operations above:
1. `Filter` to keep only log entries with "response" events
2. `Batch` the filtered items into a slice of log entries
3. `Sort` the batch by the "Src" field of the log struct

Finally, the code sets up a sink to process and display the sorted results:

```go
// Set up sink to process sorted results
strm.Into(sinks.Func(func(items []log) error {
    for _, item := range items {
        fmt.Printf("%v\n", item)
    }
    return nil
}))
```
