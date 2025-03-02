## SortByMapKey: Sorting Stream Elements by Map Key

This example demonstrates how to sort a batch of map elements in an Automi stream by a specific map key. The code processes a stream of log entries and sorts them by the URL source path.

### Code Overview
The example creates a stream of log entries (maps), filters for response events, batches them together, and sorts them by the "Src" key:

```go
func main() {
    // Define log type and sample data
    type log map[string]string
    source := sources.Slice([]log{
        {"Event": "request", "Src": "/i/a", "Device": "00:11:51:AA", "Result": "accepted"},
        {"Event": "response", "Src": "/i/a/", "Device": "00:11:51:AA", "Result": "served"},
        // More log entries...
    })

    // Create and configure stream
    strm := stream.From(source)
    strm.WithLogSink(sinks.SlogJSON(slog.LevelDebug))
    strm.Run(
        // Filter for response events only
        exec.Filter(func(_ context.Context, e log) bool {
            return (e["Event"] == "response")
        }),
        
        // Convert to map[string]string type to support downstrem map-specific operations
        exec.Map(func(_ context.Context, e log) map[string]string {
            return e
        }),
        
        // Batch items into a slice
        window.Batch[map[string]string](),
        
        // Sort by the "Src" key
        exec.SortByMapKey[[]map[string]string]("Src"),
    )

    // Print the sorted results
    strm.Into(sinks.Func(func(sorted []map[string]string) error {
        for _, item := range sorted {
            fmt.Printf("%v\n", item)
        }
        return nil
    }))
}
```

### Key Points
* The `SortByMapKey` operator sorts a batch of map elements by a specific key
* `window.Batch` collects stream elements into a slice before sorting
* Type parameters must match the data structure being processed (a slice of maps)
* This pattern is useful for organizing and presenting log data by specific attributes.

The output will look something like this:

```
map[Device:00:11:51:AA Event:response Result:served Src:/i/a/]
map[Device:00:11:22:33 Event:response Result:served Src:/i/b]
map[Device:00:11:51:AA Event:response Result:served Src:/i/c]
map[Device:00:BB:22:DD Event:response Result:served Src:/i/d]
```