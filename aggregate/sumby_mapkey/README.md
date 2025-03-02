Collecting workspace information```md
# Aggregate `SumByMapKey` Example

This example demonstrates how to use Automi to aggregate data from a slice of maps, calculate the sum of values associated with a specific key in those maps, and display the result.

## Code walkthrough

First, the code defines a slice of maps, where each map represents data with string keys and string values:

```go
data := []map[string]string{
    {"Zip": "10452", "Count": "17", "Female": "12", "Male": "5"},
    {"Zip": "10453", "Count": "14", "Female": "7", "Male": "7"},
    {"Zip": "10454", "Count": "18", "Female": "8", "Male": "10"},
    {"Zip": "10455", "Count": "27", "Female": "17", "Male": "10"},
    {"Zip": "10456", "Count": "5", "Female": "3", "Male": "2"},
    {"Zip": "10458", "Count": "52", "Female": "25", "Male": "27"},
    {"Zip": "10459", "Count": "7", "Female": "5", "Male": "2"},
    {"Zip": "10460", "Count": "27", "Female": "20", "Male": "7"},
    {"Zip": "10461", "Count": "49", "Female": "26", "Male": "23"},
}
```

Next, it creates a stream from this source and chains operations using the `Run()` method:

```go
strm := stream.From(sources.Slice(data)).Run(
    // Reduce each incoming map to a summary map
    exec.Map(func(_ context.Context, data map[string]string) map[string]int {
        count, _ := strconv.Atoi(data["Count"])
        female, _ := strconv.Atoi(data["Female"])
        male, _ := strconv.Atoi(data["Male"])
        return map[string]int{"Count": count, "Female": female, "Male": male}
    }),

    // batch all data into a single window
    window.Batch[map[string]int](),

    // sums the "Female" key from the map
    exec.SumByMapKey[[]map[string]int]("Female"),
)
```

The stream operations above:
- Converts the string values in the initial map to integers using [`strconv.Atoi`](https://pkg.go.dev/strconv#Atoi) and creates a new map with integer values.
- Accumulates all the data into a single batch using `window.Batch`.
- Sums the values associated with the "Female" key across all maps in the batch using `exec.SumByMapKey`.

Finally, the code defines a sink that prints the total number of females:

```go
strm.Into(sinks.Func(func(items float64) error {
    fmt.Printf("Total female is %.0f\n", items)
    return nil
}))
```