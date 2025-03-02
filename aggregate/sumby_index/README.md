# Aggregate Sum by Index

This example demonstrates how to use Automi to sum values at a specific index position from a collection of slices.
The code processes data from a CSV file and sums up all values at index 1 (which represents the "female" count) from the processed arrays.

## Code walkthrough

First, the code opens a CSV file and creates a source:

```go
// Open data source file
file, err := os.Open("stats.csv")
if err != nil {
    fmt.Println(err)
    return
}
defer file.Close()

// Define an Automi CSV source from file
src := sources.CSV(file).
    CommentChar('#'). // sets comment charcter, default is '#'
    DelimChar(',').   // sets delimiter char, default is  ','
    HasHeaders()      // indicate CSV has headers
```

Next, it creates a stream and defines a processing pipeline:

```go
// Create a stream from source
strm := stream.From(src)

// Define stream operations
strm.Run(
    // Reduces and maps each row to []int containing summary data
    exec.Map(func(_ context.Context, data []string) []int {
        count, _ := strconv.Atoi(data[1])
        female, _ := strconv.Atoi(data[2])
        male, _ := strconv.Atoi(data[4])
        return []int{count, female, male}
    }),

    // Batch incoming []int into a window
    window.Batch[[]int](),

    // Sum all values at pos(1)
    exec.SumByIndex[[][]int](1),
)
```

The stream operations:
1. Map each CSV row to a slice of integers `[]int` containing count, female, and male values
2. Batch those slices into a collection using `window.Batch[[]int]()`
3. Use `exec.SumByIndex[[][]int](1)` to sum all values at index position 1 (female counts)

Finally, the result is collected in a sink function which prints the total:

```go
strm.Into(sinks.Func(func(val float64) error {
    fmt.Printf("Total female is %.0f\n", val)
    return nil
}))

if err := <-strm.Open(context.Background()); err != nil {
    fmt.Println(err)
    return
}
```