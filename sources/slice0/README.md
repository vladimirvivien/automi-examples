## Slice Source

This example shows how to stream data from a Go slice source.
The code creates an Automi source using a slice as its data source:

```go
// Define slice source
slice := sources.Slice([]rune(`B世!ぽ@opqDQRS#$%^&*()...ef7ghijCklrAstvw`))
```

The remainder of the code defines operations to map, filter, batch, and sort each item from
the slice and finally print the result to stdout:

```go
	// create stream with emitter of rune slice
	strm := stream.From(slice)

	strm.Run(
		exec.Filter(func(_ context.Context, item rune) bool {
			return item >= 65 && item < (65+26) // remove unwanted chars
		}),
        
		exec.Map(func(_ context.Context, item rune) string {
			return string(item) // map rune to string
		}),

		// batch incoming string items
		window.Batch[string](),

		// sort batched items
		exec.SortSlice[[]string](),
	)

	// Send result to stdout
	strm.Into(sinks.Writer[string](os.Stdout)) 
```