# Aggregate Sum by Struct Field

This example demonstrates how to use Automi to sum values of a specific field from a collection of struct objects.
The code sums up all the `Male` field values from a slice of custom `stat` struct type.

## Code walkthrough

First, the code defines a custom `stat` struct type and creates a slice source:

```go
type stat struct {
    Zip                 string
    Count, Female, Male int
}

data := []stat{
    {Zip: "10452", Count: 17, Female: 12, Male: 5},
    {Zip: "10453", Count: 14, Female: 7, Male: 7},
    // more stat entries...
}
```

Next, it creates a stream from this source and chains operations using the `Run()` method:

```go
strm := stream.From(sources.Slice(data)).Run(
    // Batch the data into a slice of stat
    window.Batch[stat](),
    // Sum by the specified struct field
    exec.SumByStructField[[]stat]("Male"),
)
```

The stream operations above:
1. Create a source from the slice of stat objects
2. Batch the individual stat objects into a slice using `window.Batch[stat]()`
3. Use `exec.SumByStructField[[]stat]("Male")` to sum the values in the "Male" field across all stats

Finally, the result is collected in a sink function which prints the result:

```go
// Setup sink to collect sum
strm.Into(sinks.Func(func(val float64) error {
    fmt.Printf("Total male is %.0f\n", val)
    return nil
}))
```