# SortByIndex Example

This example demonstrates how to sort a batch of slice data by a specific index position using Automi. The code processes a collection of log entries (represented as string slices), filters specific entries, batches them, and sorts by a specified index position.

## Code walkthrough

First, the code creates a slice source containing log entries as string slices:

```go
src := sources.Slice([][]string{
    {"request", "/i/a", "00:11:51:AA", "accepted"},
    {"response", "/i/a/", "00:11:51:AA", "failed"},
    // more log entries...
})
```

Each slice follows the format: `[event, path, device, result]`

Next, it creates a stream from this source and chains operations using the `Run()` method:

```go
strm := stream.From(src).Run(
    // Filter out just responses
    exec.Filter(func(ctx context.Context, e []string) bool {
        return (e[0] == "response")
    }),

    // Batch items into -> [][]string
    window.Batch[[]string](),

    // Sort batched items by position index 3 (result field)
    exec.SortSliceByIndex[[][]string](3),
)
```

The stream operations:
1. `Filter` to keep only log entries where the first element (index 0) is "response"
2. `Batch` the filtered items into a slice of string slices
3. `Sort` the batch by the fourth element (index 3) of each slice, which represents the result field

Finally, the code sets up a sink to process and display the sorted results:

```go
// Set up sink to process sorted results
strm.Into(sinks.Func(func(items [][]string) error {
    for _, item := range items {
        fmt.Printf("%v\n", item)
    }
    return nil
}))
```