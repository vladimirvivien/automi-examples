## Aggregate - Sort Slice

This example shows how to use Automi to sort streamed slice items. First, the code
sets up the data source a Go channel that emits strings of comma-separated numbers.

```go
	ch := make(chan string)
	go func() {
		defer close(ch)
		ch <- "10452,17,12,0.71,5,0.29,0,0,17,100"
    ...
    }
```

Next, the code creates the stream using the channel as its source. The stream
declares two `Map` operators. The first operator transforms each strin item into a slice of strings.
The second `Map` operator emits each item as a slice of float64 `[]float64`. Finally, each streamed
slice is sorted before being displayed on stdout.

```go
	// Create a new stream to source from the channel
	strm := stream.From(sources.Chan(ch))

	// Declare stream operators 
	strm.Run(
		// Define a map to map string to []string
		exec.Map(func(_ context.Context, row string) []string {
			return strings.Split(row, ",")
		}),

		// Define a map to transform []string to []float64
		exec.Map(func(_ context.Context, data []string) []float64 {
			result := make([]float64, len(data))
			for i, d := range data {
				f, _ := strconv.ParseFloat(d, 64)
				result[i] = f
			}
			return result
		}),

		// Sort each float64 slice
		exec.SortSlice[[]float64](),
	)

	// Send stream items to stdout
	strm.Into(sinks.Writer[string](os.Stdout))
```