# Aggregate Sum Example

This example demonstrates how to use Automi to stream numeric data from a channel source and compute a sum of values. It shows how to extract specific fields from comma-separated data, convert them to numeric values, batch them, and then apply the sum operation.

## Code walkthrough

First, the code creates a Go channel and populates it with strings of comma-separated values:

```go
ch := make(chan string)
go func() {
    defer close(ch)
    ch <- "10452,17,12,0.71,5,0.29,0,0,17,100"
    ch <- "10453,14,7,0.5,7,0.5,0,0,14,100"
    // more data...
}()
```

Next, it creates a stream from this channel source:

```go
strm := stream.From(sources.Chan(ch))
```

The stream processing pipeline is defined using the `Run()` method with several operations:

```go
strm.Run(
    // Split each row into fields
    exec.Map(func(_ context.Context, row string) []string {
        return strings.Split(row, ",")
    }),

    // Extract the 4th field and convert to float
    exec.Map(func(_ context.Context, data []string) float64 {
        f, _ := strconv.ParseFloat(data[3], 32)
        return f
    }),

    // Batch incoming float64 --> []float64
    window.Batch[float64](),

    // Sum all float64 values into --> a single float64
    exec.Sum[[]float64, float64](),
)
```

The operations above:
1. Parse each row into a slice of string values
2. Extract the 4th field (index 3) and convert it to a float64
3. Batch all the float64 values into a single slice
4. Apply the `Sum` operator to compute the total

Finally, the result is sent to a sink that prints the total:

```go
strm.Into(sinks.Func(func(total float64) error {
    fmt.Printf("Total is %.2f\n", total)
    return nil
}))
```

